//
//  CMMiscPropertyProviders.m
//
//
//  Created by Steve Cosman on 2023-07-07.
//

#import "CMMiscPropertyProviders.h"

#import "../utils/CMUtils.h"

@import WatchConnectivity;

@implementation CMAppInstallDatePropertyProviders

- (NSDate *)dateValue {
    NSURL *docsFolderUrl = [[[NSFileManager defaultManager] URLsForDirectory:NSDocumentDirectory
                                                                   inDomains:NSUserDomainMask] lastObject];
    if (!docsFolderUrl) {
        return 0;
    }

    NSError *err;
    NSDate *appInstallDate = [[[NSFileManager defaultManager] attributesOfItemAtPath:docsFolderUrl.path error:&err]
        objectForKey:NSFileCreationDate];

    if (err != nil || appInstallDate == nil) {
        return nil;
    }

    return appInstallDate;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeTime;
}

@end

@interface CMHasWatchPropertyProviders () <WCSessionDelegate>

@property(nonatomic) NSNumber *hasWatch;
@property(nonatomic, strong) dispatch_group_t waitGroup;

@end

@implementation CMHasWatchPropertyProviders

- (BOOL)boolValue {
    if (!WCSession.isSupported) {
        return NO;
    }

    // Cache population
    if (!self.hasWatch) {
        @synchronized(self) {
            if (!self.waitGroup) {
                self.waitGroup = dispatch_group_create();

                WCSession *defaultSession = WCSession.defaultSession;
                defaultSession.delegate = self;
                [defaultSession activateSession];

                dispatch_group_enter(self.waitGroup);
            }
        }

        dispatch_group_wait(self.waitGroup, dispatch_time(DISPATCH_TIME_NOW, 2.0 * NSEC_PER_SEC));
    }

    if (!self.hasWatch) {
        return NO;
    }
    return self.hasWatch.boolValue;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

#pragma Mark WCSessionDelegate

- (void)session:(nonnull WCSession *)session
    activationDidCompleteWithState:(WCSessionActivationState)activationState
                             error:(nullable NSError *)error {
    if (!error && activationState == WCSessionActivationStateActivated) {
        self.hasWatch = [NSNumber numberWithBool:session.paired];
    }

    dispatch_group_leave(self.waitGroup);
}

- (void)sessionDidBecomeInactive:(nonnull WCSession *)session {
}

- (void)sessionDidDeactivate:(nonnull WCSession *)session {
}

@end
