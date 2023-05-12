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
    // TODO: move to bg thread?

    // Register the action dispatcher
    [CMLibBindings registerWithAppcore];

    NSError *error;
    [AppcoreSharedAppcore() start:&error];
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
