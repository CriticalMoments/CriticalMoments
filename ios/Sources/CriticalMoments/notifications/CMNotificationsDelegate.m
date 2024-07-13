//
//  CMNotificationsDelegate.m
//
//
//  Created by Steve Cosman on 2024-07-06.
//

#import "CMNotificationsDelegate.h"

#import "../CriticalMoments_private.h"

@import Appcore;

@interface CMNotificationsDelegate ()

@property(nonatomic, weak) id<UNUserNotificationCenterDelegate> originalDelegate;

@end

@implementation CMNotificationsDelegate

- (id)initWithOriginalDelegate:(id<UNUserNotificationCenterDelegate>)originalDelegate {
    self = [super init];
    if (self) {
        self.originalDelegate = originalDelegate;
    }
    return self;
}

- (void)userNotificationCenter:(UNUserNotificationCenter *)center
       willPresentNotification:(UNNotification *)notification
         withCompletionHandler:(void (^)(UNNotificationPresentationOptions))completionHandler {
    // Prefer original delegate's behaviour if available
    if ([self.originalDelegate respondsToSelector:@selector(userNotificationCenter:
                                                           willPresentNotification:withCompletionHandler:)]) {
        [self.originalDelegate userNotificationCenter:center
                              willPresentNotification:notification
                                withCompletionHandler:completionHandler];
        return;
    }

    // Align to the default OS behaviour if the app doesn't set a delegate: don't display any notifications overtop the
    // active app.
    completionHandler(0);
}

- (void)userNotificationCenter:(UNUserNotificationCenter *)center
    didReceiveNotificationResponse:(UNNotificationResponse *)response
             withCompletionHandler:(void (^)(void))completionHandler {

    // Handle this notification if it's for CriticalMoments
    NSString *notifId = response.notification.request.identifier;
    if ([notifId hasPrefix:DatamodelNotificationUniqueIDPrefix]) {
        [CriticalMoments.shared actionForNotification:notifId];
        completionHandler();
        return;
    }

    if ([self.originalDelegate respondsToSelector:@selector(userNotificationCenter:
                                                      didReceiveNotificationResponse:withCompletionHandler:)]) {
        [self.originalDelegate userNotificationCenter:center
                       didReceiveNotificationResponse:response
                                withCompletionHandler:completionHandler];
        return;
    }

    completionHandler();
}

- (void)userNotificationCenter:(UNUserNotificationCenter *)center
    openSettingsForNotification:(UNNotification *)notification {
    [self.originalDelegate userNotificationCenter:center openSettingsForNotification:notification];
}

#ifdef DEBUG
// Only compiled in debug mode, not included in release builds.
// Check everything is setup correctly, and log a warning if not.
+ (void)devModeCheckNotificationDelegateSetupCorrectly {
    // Check the app's notification delegate is ours (potentually wrapping theirs)
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    id<UNUserNotificationCenterDelegate> currentDelegate = center.delegate;
    if (!currentDelegate || ![currentDelegate isKindOfClass:[CMNotificationsDelegate class]]) {
        NSLog(@"CriticalMoments Debug Developer Warning: The CM notification delegate is not registered. Tapping "
              @"notifications from CM won't trigger the correct actions. If you register a custom "
              @"UNUserNotificationCenterDelegate, do so before calling the start method of CriticalMoments.\nThis "
              @"warning is only in debug builds, and is not included in release builds.");
    }
}
#endif

@end
