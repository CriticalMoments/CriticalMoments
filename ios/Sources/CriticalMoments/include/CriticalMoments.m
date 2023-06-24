//
//  CriticalMoments.m
//  CriticalMoments
//
//  Created by Steve Cosman on 2023-04-17.
//

#import "CriticalMoments.h"

#import "../appcore_integration/CMLibBindings.h"
#import "../properties/CMPropertyRegisterer.h"

@implementation CriticalMoments

+ (NSString *)objcPing {
    return @"objcPong";
}

+ (NSString *)goPing {
    return AppcoreGoPing();
}

+ (void)start {
    // Nested dispatch to main then background. Why?
    // We want critical moments to start on background thread, but we want it to
    // start after the app setup is done. Some property providers will provide
    // unknown values before the main thread is ready. This puts CM startup
    // after core app setup.
    dispatch_async(dispatch_get_main_queue(), ^{
      dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        NSError *error = [CriticalMoments startReturningError];
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

+ (NSError *)startReturningError {
    // Register the action dispatcher and properties
    [CMLibBindings registerWithAppcore];

    CMPropertyRegisterer *propertryRegisterer = [[CMPropertyRegisterer alloc] init];
    [propertryRegisterer registerDefaultPropertiesToAppcore];

    // Set the cache directory to applicationSupport/CriticalMomentsData
    NSURL *appSupportDir = [[NSFileManager.defaultManager URLsForDirectory:NSApplicationSupportDirectory
                                                                 inDomains:NSUserDomainMask] lastObject];
    NSURL *criticalMomentsCacheDir = [appSupportDir URLByAppendingPathComponent:@"CriticalMomentsData"];
    NSError *error;
    [NSFileManager.defaultManager createDirectoryAtURL:criticalMomentsCacheDir
                           withIntermediateDirectories:YES
                                            attributes:nil
                                                 error:&error];
    if (error) {
        return error;
    }
    [AppcoreSharedAppcore() setCacheDirPath:[criticalMomentsCacheDir path] error:&error];
    if (error) {
        return error;
    }

    [AppcoreSharedAppcore() start:&error];
    if (error) {
        return error;
    }
    return nil;
}

+ (void)setConfigUrl:(NSString *)urlString {
    NSError *error;
    [AppcoreSharedAppcore() setConfigUrl:urlString error:&error];
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

+ (void)sendEvent:(NSString *)eventName {
    NSError *error;
    [AppcoreSharedAppcore() sendEvent:eventName error:&error];
    if (error) {
        NSLog(@"WARN: CriticalMoments -- error sending event: %@", error);
    }
}

@end
