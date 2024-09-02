//
//  CriticalMoments.m
//  CriticalMoments
//
//  Created by Steve Cosman on 2023-04-17.
//

#import "CriticalMoments.h"

#import "../CriticalMoments_private.h"
#import "../appcore_integration/CMLibBindings.h"
#import "../background/CMBackgroundHandler.h"
#import "../messaging/CMBannerManager.h"
#import "../notifications/CMNotificationHandler.h"
#import "../notifications/CMNotificationsDelegate.h"
#import "../properties/CMPropertyRegisterer.h"
#import "../themes/CMTheme_private.h"
#import "../utils/CMNotificationObserver.h"
#import "../utils/CMUtils.h"

#import <UserNotifications/UserNotifications.h>
#import <os/log.h>

@interface CriticalMoments ()
@property(nonatomic) BOOL queuesStarted, disableNotifications;
@property(nonatomic, strong) NSString *releaseConfigUrl, *devConfigUrl;
@property(nonatomic, strong) AppcoreAppcore *appcore;
@property(nonatomic, strong) CMLibBindings *bindings;
@property(nonatomic, strong) CMTheme *currentTheme;
@property(nonatomic, strong) CMNotificationObserver *notificationObserver;
@property(nonatomic, strong) dispatch_queue_t actionQueue;
@property(nonatomic, strong) dispatch_queue_t eventQueue;
@property(nonatomic, strong) CMNotificationsDelegate *notificationDelegate;
@property(nonatomic, strong) CMNotificationHandler *notificationHandler;
@end

@implementation CriticalMoments

- (id)initInternal {
    self = [super init];
    if (self) {
        _appcore = AppcoreNewAppcore();
        _notificationObserver = [[CMNotificationObserver alloc] initWithCm:self];

        // Event queue is serial to preserve event order
        dispatch_queue_attr_t eventQueueAttr =
            dispatch_queue_attr_make_with_qos_class(DISPATCH_QUEUE_SERIAL, QOS_CLASS_DEFAULT, 0);
        _eventQueue = dispatch_queue_create("io.criticalmoments.event_queue", eventQueueAttr);

        // action queue is concurrent
        dispatch_queue_attr_t actionQueueAttr =
            dispatch_queue_attr_make_with_qos_class(DISPATCH_QUEUE_CONCURRENT, QOS_CLASS_DEFAULT, 0);
        _actionQueue = dispatch_queue_create("io.criticalmoments.action_queue", actionQueueAttr);

        // action and event queues suspended until we start (not just called start,
        // but fully started)
        dispatch_suspend(_eventQueue);
        dispatch_suspend(_actionQueue);
    }
    return self;
}

- (void)dealloc {
    // Can't deallocate suspended queues. Typically no-op but in tests is needed.
    [self startQueues];
}

static CriticalMoments *sharedInstance = nil;

+ (CriticalMoments *)shared {
    // avoid lock if we can
    if (sharedInstance) {
        return sharedInstance;
    }

    @synchronized(CriticalMoments.class) {
        if (!sharedInstance) {
            sharedInstance = [[self alloc] initInternal];
        }

        return sharedInstance;
    }
}

+ (CriticalMoments *)sharedInstance {
    return [CriticalMoments shared];
}

- (NSString *)objcPing {
    return @"objcPong";
}

- (NSString *)goPing {
    return AppcoreGoPing();
}

- (AppcoreAppcore *)appcore {
    return _appcore;
}

- (void)start {
    // Start notification observer before main queues, so that the enter_forground and other events are at head of queue
    [self.notificationObserver start];

    // Register notification delegate on main thread. This must be done before app finishes launching, so can't be
    // deferred.
    [self registerNotificationDelegate];

    // Registering BG work must be done before the end of app launch, do no defer
    //  https://developer.apple.com/documentation/backgroundtasks/bgtaskscheduler/register(fortaskwithidentifier:using:launchhandler:)?language=objc
    [self registerBackgroundHandler];

    // Nested dispatch to main then background. Why?
    // We want critical moments to start on background thread, but we want it to
    // start after the app setup is done. Some property providers will provide
    // unknown values before the main thread is ready. This puts CM startup
    // after core app setup.
    dispatch_async(dispatch_get_main_queue(), ^{
      dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        NSError *error = [self startReturningError];
        if (error) {
            os_log_fault(OS_LOG_DEFAULT, "CriticalMoments: Critical Moments was unable to start!\nCMError: %@",
                         error.localizedDescription);
#if DEBUG
            os_log_fault(OS_LOG_DEFAULT,
                         "CriticalMoments: throwing a NSInternalInconsistencyException to help find the issue above. "
                         "Exceptions are only thrown in debug builds. This will not crash apps built for release.");
            @throw NSInternalInconsistencyException;
#endif
        }
      });
    });
}

- (NSError *)startReturningError {
#if DEBUG
    if (!_releaseConfigUrl) {
        // Warn the developer if they have not set a release URL
        NSLog(@"CriticalMoments: you have not set a valid release config url. Critical Moments will not work in "
              @"release builds/app-store builds. Be sure to set this before releasing your app.");
    }
    if (_devConfigUrl) {
        [self setConfigUrl:_devConfigUrl];
    } else if (_releaseConfigUrl) {
        // devConfig is optional. Fall back to release if that's all they set
        [self setConfigUrl:_releaseConfigUrl];
    }
#else
    // Only use the production URL on non debug builds
    if (_releaseConfigUrl) {
        [self setConfigUrl:_releaseConfigUrl];
    }
#endif

    // Register the action dispatcher and properties
    if (!self.bindings) {
        self.bindings = [[CMLibBindings alloc] initWithCM:self];
    }
    [_appcore registerLibraryBindings:_bindings];

    // Fix the timezone -- golang doesn't know local offset by default
    NSTimeZone *tz = NSTimeZone.localTimeZone;
    [_appcore setTimezoneGMTOffset:tz.secondsFromGMT];

    CMPropertyRegisterer *propertryRegisterer = [[CMPropertyRegisterer alloc] initWithAppcore:_appcore];
    [propertryRegisterer registerDefaultPropertiesToAppcore];

    NSURL *appSupportDir = [[NSFileManager.defaultManager URLsForDirectory:NSApplicationSupportDirectory
                                                                 inDomains:NSUserDomainMask] lastObject];

    // Set the data directory to applicationSupport/critical_moments_data
    NSError *error;
    NSURL *criticalMomentsDataDir = [appSupportDir URLByAppendingPathComponent:@"critical_moments_data"];

    [NSFileManager.defaultManager createDirectoryAtURL:criticalMomentsDataDir
                           withIntermediateDirectories:YES
                                            attributes:nil
                                                 error:&error];

    if (error) {
        return error;
    }
    [_appcore setDataDirPath:[criticalMomentsDataDir path] error:&error];
    if (error) {
        return error;
    }

    bool allowDebugLoad = [self isDebug] && [self urlAllowedForDebugLoad];
    [_appcore start:allowDebugLoad error:&error];
    if (error) {
        return error;
    }

    // We've started now. Can resume the two worker queues.
    [self startQueues];

#ifdef DEBUG
    // check everything is setup correctly
    [self devModeCheckSetupCorrectly];
#endif

    return nil;
}

- (bool)urlAllowedForDebugLoad {
    // Allow any in test cases
    if ([@"com.apple.dt.xctest.tool" isEqualToString:NSBundle.mainBundle.bundleIdentifier]) {
        return true;
    }
    // Allowed if in main bundle, but not if in data directories (downloaded/external)
    NSString *mainBundlePath = NSBundle.mainBundle.bundleURL.absoluteString;
    if (!mainBundlePath) {
        return false;
    }
    return [_appcore.configUrl hasPrefix:mainBundlePath];
}

- (void)startQueues {
    @synchronized(self) {
        if (!_queuesStarted) {
            _queuesStarted = true;
            dispatch_resume(_eventQueue);
            dispatch_resume(_actionQueue);
        }
    }
}

- (void)setApiKey:(NSString *)apiKey error:(NSError **)returnError {
    // Set API Key
    NSError *error;
    NSString *bundleIdentifier = [[NSBundle mainBundle] bundleIdentifier];
    [_appcore setApiKey:apiKey bundleID:bundleIdentifier error:&error];
    if (error) {
        if (returnError) {
            *returnError = error;
        }
        NSLog(@"ERROR: CriticalMoments -- [Invalid API Key]: %@", error);
#if DEBUG
        if (!returnError) {
            NSLog(@"CriticalMoments: throwing a NSInternalInconsistencyException "
                  @"to help find this issue. Exceptions are only thrown in debug "
                  @"mode and when you don't pass and error to detect/handle the issue. "
                  @"This will not crash apps built for release.");
            @throw NSInternalInconsistencyException;
        }
#endif
    }
}

- (nonnull NSString *)getApiKey {
    NSString *apiKey = [_appcore apiKey];
    if (apiKey.length == 0) {
        return nil;
    }
    return apiKey;
}

- (void)setDevelopmentConfigName:(NSString *)configFileName {
    BOOL success = [self setDevelopmentConfigNameWithSuccess:configFileName fromBundle:nil];
    if (!success) {
        os_log_fault(OS_LOG_DEFAULT, "CriticalMoments: unable to find config file: %@", configFileName);
    }
}

- (BOOL)setDevelopmentConfigNameWithSuccess:(NSString *)configFileName fromBundle:(NSBundle *_Nullable)bundle {
    if (!bundle) {
        bundle = NSBundle.mainBundle;
    }
    NSString *extension = [configFileName pathExtension];
    NSString *resourceName = [configFileName stringByDeletingPathExtension];
    NSURL *localConfigUrl = [bundle URLForResource:resourceName withExtension:extension];
    if (!localConfigUrl) {
        return false;
    }
    NSString *filePath = [NSString stringWithFormat:@"file://%@", localConfigUrl.path];
    [self setDevelopmentConfigUrl:filePath];
    return true;
}

- (void)setDevelopmentConfigUrl:(NSString *)urlString {
    if (![urlString hasPrefix:@"file://"]) {
        NSLog(@"CriticalMoments: invalid file URL sent to setDevelopmentConfigUrl. The URL must begin with `file://`");
        return;
    }

    _devConfigUrl = urlString;
}

- (void)setReleaseConfigUrl:(NSString *)urlString {
    if (![urlString hasPrefix:@"https://"]) {
        NSLog(@"CriticalMoments: invalid URL sent to setProductionConfigUrl. The URL must begin with `https://`");
        return;
    }

    _releaseConfigUrl = urlString;
}

- (void)setConfigUrl:(NSString *)urlString {
    NSError *error;
    [_appcore setConfigUrl:urlString error:&error];
    if (error != nil) {
        NSLog(@"ERROR: CriticalMoments -- invalid remote config url: %@", error);
#if DEBUG
        NSLog(@"CriticalMoments: throwing a NSInternalInconsistencyException "
              @"to help find this issue. Exceptions are only thrown in debug "
              @"mode, and will not crash apps built for release.");
        @throw NSInternalInconsistencyException;
#endif
    }
}

- (void)setLogEvents:(bool)logEvents {
    [self.appcore setLogEvents:logEvents];
}

- (void)sendEvent:(NSString *)eventName {
    [self sendEvent:eventName builtIn:false handler:nil];
}

- (void)sendEvent:(NSString *)eventName builtIn:(bool)builtIn handler:(void (^)(NSError *_Nullable error))handler {
    __block void (^blockHandler)(NSError *_Nullable error) = handler;
    __block NSString *blockEventName = eventName;
    dispatch_async(_eventQueue, ^{
      NSError *error;
      if (builtIn) {
          [_appcore sendBuiltInEvent:blockEventName error:&error];
      } else {
          [_appcore sendClientEvent:blockEventName error:&error];
      }
      if (error) {
          NSLog(@"CriticalMoments: Error sending event. %@", error.localizedDescription);
      }

      if (blockHandler) {
          blockHandler(error);
      }
    });
}

- (void)checkNamedCondition:(NSString *)name handler:(void (^)(bool, NSError *_Nullable))handler {
    [self checkNamedCondition:name completionHandler:handler];
}

- (void)checkNamedCondition:(NSString *)name completionHandler:(void (^_Nonnull)(bool, NSError *_Nullable))result {
    __block void (^blockHandler)(bool result, NSError *_Nullable error) = result;
    __block NSString *blockName = name;
    dispatch_async(_actionQueue, ^{
      NSError *error;
      BOOL result;
      [_appcore checkNamedCondition:blockName returnResult:&result error:&error];

      if (blockHandler) {
          blockHandler(result, error);
      }
    });
}

// Private, only for internal use (demo app).
// Do not use this in any other apps. Will always return error.
- (void)checkInternalTestCondition:(NSString *_Nonnull)conditionString
                           handler:(void (^_Nonnull)(bool result, NSError *_Nullable error))handler {
    __block void (^blockHandler)(bool result, NSError *_Nullable error) = handler;
    __block NSString *blockCondition = conditionString;
    dispatch_async(_actionQueue, ^{
      NSError *error;
      BOOL result;
      [_appcore checkTestCondition:blockCondition returnResult:&result error:&error];
      if (error) {
          NSLog(@"CriticalMoments: error in test func %@", error);
      }

      if (blockHandler) {
          blockHandler(result, error);
      }
    });
}

- (void)performNamedAction:(NSString *)name handler:(void (^)(NSError *_Nullable))handler {
    __block void (^blockHandler)(NSError *_Nullable error) = handler;
    __block NSString *blockName = name;
    dispatch_async(_actionQueue, ^{
      NSError *error;
      [_appcore performNamedAction:blockName error:&error];
      if (blockHandler) {
          blockHandler(error);
      }
    });
}

- (DatamodelTheme *)themeFromConfigByName:(NSString *)name {
    return [_appcore themeForName:name];
}

#pragma mark Custom Properties

- (void)registerStringProperty:(NSString *)value forKey:(NSString *)name error:(NSError *_Nullable *)error {
    [self setStringProperty:value forKey:name error:error];
}

- (BOOL)setStringProperty:(NSString *)value forKey:(NSString *)name error:(NSError *_Nullable *)error {
    return [_appcore registerClientStringProperty:name value:value error:error];
}

- (void)registerBoolProperty:(BOOL)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error {
    [self setBoolProperty:value forKey:name error:error];
}

- (BOOL)setBoolProperty:(BOOL)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error {
    return [_appcore registerClientBoolProperty:name value:value error:error];
}

- (void)registerFloatProperty:(double)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error {
    [self setFloatProperty:value forKey:name error:error];
}

- (BOOL)setFloatProperty:(double)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error {
    return [_appcore registerClientFloatProperty:name value:value error:error];
}

- (void)registerIntegerProperty:(long long)value
                         forKey:(NSString *)name
                          error:(NSError *_Nullable __autoreleasing *)error {
    [self setIntegerProperty:value forKey:name error:error];
}

- (BOOL)setIntegerProperty:(long long)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error {
    return [_appcore registerClientIntProperty:name value:value error:error];
}

- (void)registerTimeProperty:(NSDate *)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error {
    [self setTimeProperty:value forKey:name error:error];
}

- (BOOL)setTimeProperty:(NSDate *)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error {
    int64_t goTime = [CMUtils dateToGoTime:value];
    return [_appcore registerClientTimeProperty:name value:goTime error:error];
}

- (void)registerPropertiesFromJson:(NSData *)jsonData error:(NSError *_Nullable __autoreleasing *)error {
    [self setPropertiesFromJson:jsonData error:error];
}

- (BOOL)setPropertiesFromJson:(NSData *)jsonData error:(NSError *_Nullable __autoreleasing *)error {
    return [_appcore registerClientPropertiesFromJson:jsonData error:error];
}

#pragma mark Current Theme

- (CMTheme *)currentTheme {
    // avoid lock if we can
    if (!_currentTheme) {
        @synchronized(self) {
            if (!_currentTheme) {
                _currentTheme = [[CMTheme alloc] init];
            }
        }
    }
    return _currentTheme;
}

#pragma mark Background Work

- (void)registerBackgroundHandler {
    CMBackgroundHandler *bgh = [[CMBackgroundHandler alloc] initWithCm:self];
    self.backgroundHandler = bgh;
    [bgh registerBackgroundTasks];
}

- (void)runAppcoreBackgroundWork {
    NSError *error;
    [_appcore performBackgroundWork:&error];
    if (error) {
        os_log_error(OS_LOG_DEFAULT, "CriticalMoments: issue performing background work");
    }
}

#pragma mark Notifications

- (void)disableUserNotifications {
    _disableNotifications = true;
}

- (BOOL)userNotificationsDisabled {
    return _disableNotifications;
}

- (AppcoreNotificationPlan *)currentNotificationPlan:(NSError *_Nullable *_Nullable)error {
    return [_appcore fetchNotificationPlan:error];
}

- (void)updateNotificationPlan:(AppcoreNotificationPlan *_Nullable)notifPlan {
    [self.notificationHandler updateNotificationPlan:notifPlan];
}

- (void)registerNotificationDelegate {
    if (_disableNotifications) {
        return;
    }
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    id<UNUserNotificationCenterDelegate> existingDelegate = center.delegate;
    _notificationDelegate = [[CMNotificationsDelegate alloc] initWithOriginalDelegate:existingDelegate andCm:self];
    center.delegate = _notificationDelegate;

    CMNotificationHandler *notificationHandler = [[CMNotificationHandler alloc] initWithCm:self];
    self.notificationHandler = notificationHandler;
}

- (void)actionForNotification:(NSString *)identifier {
    dispatch_async(_actionQueue, ^{
      NSError *error;
      [_appcore actionForNotification:identifier error:&error];
      if (error) {
          NSLog(@"CriticalMoments: error in notification: %@", error.localizedDescription);
      }
    });
}

- (void)requestNotificationPermissionWithCompletionHandler:
    (void (^_Nullable)(BOOL prompted, BOOL granted, NSError *__nullable error))completionHandler {
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];

    // Check if already authorized/denied.
    // We don't want to force update the plan unless permissions are actually changing.
    [center getNotificationSettingsWithCompletionHandler:^(UNNotificationSettings *_Nonnull settings) {
      BOOL authorized = settings.authorizationStatus == UNAuthorizationStatusAuthorized;
      BOOL denied = settings.authorizationStatus == UNAuthorizationStatusDenied;
      if (authorized || denied) {
          if (completionHandler) {
              completionHandler(false, authorized, nil);
          }
          return;
      }
      UNAuthorizationOptions opt = UNAuthorizationOptionAlert | UNAuthorizationOptionBadge |
                                   UNAuthorizationOptionSound | UNAuthorizationOptionCriticalAlert;
      [center requestAuthorizationWithOptions:opt
                            completionHandler:^(BOOL granted, NSError *_Nullable error) {
                              if (completionHandler) {
                                  completionHandler(true, granted, error);
                              }
                              if (granted) {
                                  // Schedule any CM notifications that need to be scheduled now
                                  [_appcore forceUpdateNotificationPlan:nil];
                              }
                            }];
    }];
}

#pragma mark Demo App and Internal APIs

// Private/Internal: Only for demo app usage
- (void)setTheme:(CMTheme *)theme {
    @synchronized(self) {
        _currentTheme = theme;
    }
}

// Private/Internal: Only for demo app usage
- (void)setBuiltInTheme:(NSString *)themeName {
    CMTheme *libTheme = [CMTheme libaryThemeByName:themeName];
    if (libTheme) {
        [self setTheme:libTheme];
        return;
    }

    DatamodelTheme *dmTheme = [self themeFromConfigByName:themeName];
    if (dmTheme) {
        CMTheme *theme = [CMTheme themeFromAppcoreTheme:dmTheme];
        [self setTheme:theme];
    }
}

// Private/Internal: Only for demo app usage
- (int)builtInBaseThemeCount {
    return (int)DatamodelBaseThemeCount();
}

- (bool)isDebug {
#if DEBUG
    return true;
#endif
    return false;
}

- (void)removeAllBanners {
    [CMBannerManager.shared removeAllAppWideMessages];
}

#ifdef DEBUG
// Not compiled into release builds, only debug.
// Checks everything is setup correctly and setup logs issues for devs to see
- (void)devModeCheckSetupCorrectly {
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(2 * NSEC_PER_SEC)),
                   dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
                     [CMBackgroundHandler devModeCheckBackgroundSetupCorrectly];
                     [CMNotificationsDelegate devModeCheckNotificationDelegateSetupCorrectly];
                   });
}
#endif

@end
