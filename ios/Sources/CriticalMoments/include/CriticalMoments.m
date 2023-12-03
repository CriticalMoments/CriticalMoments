//
//  CriticalMoments.m
//  CriticalMoments
//
//  Created by Steve Cosman on 2023-04-17.
//

#import "CriticalMoments.h"

#import "../messaging/CMBannerManager.h"

#import "../appcore_integration/CMLibBindings.h"
#import "../properties/CMPropertyRegisterer.h"

@interface CriticalMoments ()
@property(nonatomic) BOOL queuesStarted;
@property(nonatomic, strong) NSString *releaseConfigUrl, *devConfigUrl;
@property(nonatomic, strong) AppcoreAppcore *appcore;
@property(nonatomic, strong) CMLibBindings *bindings;
@property(nonatomic, strong) CMTheme *currentTheme;
@property(nonatomic, strong) dispatch_queue_t actionQueue;
@property(nonatomic, strong) dispatch_queue_t eventQueue;
@end

@implementation CriticalMoments

- (id)initInternal {
    self = [super init];
    if (self) {
        _appcore = AppcoreNewAppcore();

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

+ (CriticalMoments *)sharedInstance {
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
    // Nested dispatch to main then background. Why?
    // We want critical moments to start on background thread, but we want it to
    // start after the app setup is done. Some property providers will provide
    // unknown values before the main thread is ready. This puts CM startup
    // after core app setup.
    dispatch_async(dispatch_get_main_queue(), ^{
      dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        NSError *error = [self startReturningError];
        if (error) {
            NSLog(@"CriticalMoments: Critical Moments was unable to start! "
                  @"%@",
                  error);
#if DEBUG
            NSLog(@"CriticalMoments: throwing a "
                  @"NSInternalInconsistencyException "
                  @"to help find this issue. Exceptions are only thrown in "
                  @"debug "
                  @"mode, and will not crash apps built for release.");
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

    /*BOOL success = [NSFileManager.defaultManager removeItemAtURL:criticalMomentsDataDir error:&error];
    if (!success || error) {
        NSLog(@"error removing existing cache: %@", error);
        return error;
    }*/

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

- (void)sendEvent:(NSString *)eventName {
    [self sendEvent:eventName handler:nil];
}

- (void)sendEvent:(NSString *)eventName handler:(void (^)(NSError *_Nullable error))handler {
    __block void (^blockHandler)(NSError *_Nullable error) = handler;
    __block NSString *blockEventName = eventName;
    dispatch_async(_eventQueue, ^{
      NSError *error;
      [_appcore sendClientEvent:blockEventName error:&error];

      if (blockHandler) {
          blockHandler(error);
      }
    });
}

- (void)checkNamedCondition:(NSString *)name
                  condition:(NSString *)condition
                    handler:(void (^)(bool, NSError *_Nullable))handler {
#if DEBUG
    NSError *collisionError;
    [_appcore checkNamedConditionCollision:name conditionString:condition error:&collisionError];
    if (collisionError != nil) {
        NSLog(@"\nWARNING: CriticalMoments\nWARNING: CriticalMoments\nIssue with checkNamedCondition usage. Note: this "
              @"error log is only shown when debugger attached.\n%@\n\n",
              collisionError.localizedDescription);
    }
#endif

    __block void (^blockHandler)(bool result, NSError *_Nullable error) = handler;
    __block NSString *blockName = name;
    __block NSString *blockCondition = condition;
    dispatch_async(_actionQueue, ^{
      NSError *error;
      BOOL result;
      [_appcore checkNamedCondition:blockName conditionString:blockCondition returnResult:&result error:&error];

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
    [_appcore registerClientStringProperty:name value:value error:error];
}

- (void)registerBoolProperty:(BOOL)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error {
    [_appcore registerClientBoolProperty:name value:value error:error];
}

- (void)registerFloatProperty:(double)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error {
    [_appcore registerClientFloatProperty:name value:value error:error];
}

- (void)registerIntegerProperty:(long long)value
                         forKey:(NSString *)name
                          error:(NSError *_Nullable __autoreleasing *)error {
    [_appcore registerClientIntProperty:name value:value error:error];
}

- (void)registerTimeProperty:(NSDate *)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error {
    if (!value) {
        [_appcore registerClientTimeProperty:name value:AppcoreLibPropertyProviderNilIntValue error:error];
    } else {
        int64_t epochMilliseconds = [@(floor([value timeIntervalSince1970] * 1000)) longLongValue];
        [_appcore registerClientTimeProperty:name value:epochMilliseconds error:error];
    }
}

- (void)registerPropertiesFromJson:(NSData *)jsonData error:(NSError *_Nullable __autoreleasing *)error {
    [_appcore registerClientPropertiesFromJson:jsonData error:error];
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

- (void)setTheme:(CMTheme *)theme {
    @synchronized(self) {
        _currentTheme = theme;
    }
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

@end
