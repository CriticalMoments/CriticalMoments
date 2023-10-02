//
//  CriticalMoments.m
//  CriticalMoments
//
//  Created by Steve Cosman on 2023-04-17.
//

#import "CriticalMoments.h"

#import "../appcore_integration/CMLibBindings.h"
#import "../properties/CMPropertyRegisterer.h"

@interface CriticalMoments ()
@property(nonatomic, strong) AppcoreAppcore *appcore;
@property(nonatomic, strong) CMLibBindings *bindings;
@end

@implementation CriticalMoments

- (id)initInternal {
    self = [super init];
    if (self) {
        _appcore = AppcoreNewAppcore();
    }
    return self;
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
    // Register the action dispatcher and properties
    if (!self.bindings) {
        self.bindings = [[CMLibBindings alloc] init];
    }
    [_appcore registerLibraryBindings:_bindings];

    // Fix the timezone -- golang doesn't know local offset by default
    NSTimeZone *tz = NSTimeZone.localTimeZone;
    [_appcore setTimezoneGMTOffset:tz.secondsFromGMT];

    CMPropertyRegisterer *propertryRegisterer = [[CMPropertyRegisterer alloc] initWithAppcore:_appcore];
    [propertryRegisterer registerDefaultPropertiesToAppcore];

    // Set the cache directory to applicationSupport/CriticalMomentsData
    NSURL *appSupportDir = [[NSFileManager.defaultManager URLsForDirectory:NSApplicationSupportDirectory
                                                                 inDomains:NSUserDomainMask] lastObject];
    NSError *error;
    NSURL *criticalMomentsCacheDir = [appSupportDir URLByAppendingPathComponent:@"CriticalMomentsData"];
    [NSFileManager.defaultManager createDirectoryAtURL:criticalMomentsCacheDir
                           withIntermediateDirectories:YES
                                            attributes:nil
                                                 error:&error];
    if (error) {
        return error;
    }
    [_appcore setCacheDirPath:[criticalMomentsCacheDir path] error:&error];
    if (error) {
        return error;
    }

    [_appcore start:&error];
    if (error) {
        return error;
    }
    return nil;
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
    NSError *error;
    // TODO: check we've started. Will crash otherwise
    [_appcore sendEvent:eventName error:&error];
    if (error) {
        NSLog(@"WARN: CriticalMoments -- error sending event: %@", error);
    }
}

- (bool)checkNamedCondition:(NSString *)name condition:(NSString *)condition error:(NSError **)returnError {
    // TODO: check we've started. Will crash otherwise
#if DEBUG
    NSError *collisionError;
    bool colResult = [_appcore checkNamedConditionCollision:name conditionString:condition error:&collisionError];
    if (collisionError != nil) {
        NSLog(@"\nWARNING: CriticalMoments\nWARNING: CriticalMoments\nIssue with checkNamedCondition usage. Note: this "
              @"error log is only shown when debugger attached.\n%@\n\n",
              collisionError.localizedDescription);
    }
#endif

    NSError *error;
    BOOL result;
    [_appcore checkNamedCondition:name conditionString:condition ret0_:&result error:returnError];

    if (returnError) {
        NSLog(@"ERROR: CriticalMoments -- error in checkNamedCondition: %@", (*returnError).localizedDescription);
    }

    return result;
}

- (void)performNamedAction:(NSString *)name error:(NSError **)error {
    [_appcore performNamedAction:name error:error];
}

- (DatamodelTheme *)themeFromConfigByName:(NSString *)name {
    return [_appcore themeForName:name];
}

@end
