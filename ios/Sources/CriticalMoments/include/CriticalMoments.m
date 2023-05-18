//
//  CriticalMoments.m
//  CriticalMoments
//
//  Created by Steve Cosman on 2023-04-17.
//

#import "CriticalMoments.h"

#import "../appcore_integration/CMLibBindings.h"

@import Appcore;

@implementation CriticalMoments

+ (NSString *)objcPing {
    return @"objcPong";
}

+ (NSString *)goPing {
    return AppcoreGoPing();
}

+ (void)start {
    NSError *error = [CriticalMoments startReturningError];
    if (error) {
        NSLog(@"CriticalMoments: Critical Moments was unable to start! %@",
              error);
#if DEBUG
        NSLog(@"CriticalMoments: throwing a NSInternalInconsistencyException "
              @"to help find this issue. Exceptions are only thrown in debug "
              @"mode, and will not crash apps built for release.");
        @throw NSInternalInconsistencyException;
#endif
    }
}

+ (NSError *)startReturningError {
    // TODO: move to bg thread?

    // Register the action dispatcher
    [CMLibBindings registerWithAppcore];

    // Set the cache directory to applicationSupport/CriticalMomentsData
    NSURL *appSupportDir = [[NSFileManager.defaultManager
        URLsForDirectory:NSApplicationSupportDirectory
               inDomains:NSUserDomainMask] lastObject];
    NSURL *criticalMomentsCacheDir =
        [appSupportDir URLByAppendingPathComponent:@"CriticalMomentsData"];
    NSError *error;
    [NSFileManager.defaultManager createDirectoryAtURL:criticalMomentsCacheDir
                           withIntermediateDirectories:YES
                                            attributes:nil
                                                 error:&error];
    if (error) {
        return error;
    }
    [AppcoreSharedAppcore() setCacheDirPath:[criticalMomentsCacheDir path]
                                      error:&error];
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
        NSLog(@"ERROR: CriticalMoments -- invalid remote config url: %@",
              error);
#if DEBUG
        NSLog(@"CriticalMoments: throwing a NSInternalInconsistencyException "
              @"to help find this issue. Exceptions are only thrown in debug "
              @"mode, and will not crash apps built for release.");
        @throw NSInternalInconsistencyException;
#endif
    }
}

+ (void)sendEvent:(NSString *)eventName {
    [AppcoreSharedAppcore() sendEvent:eventName];
}

@end
