//
//  CMNotificationsDelegate.m
//
//
//  Created by Steve Cosman on 2024-07-06.
//

#import "CMNotificationsDelegate.h"

#import "../CriticalMoments_private.h"
#import <os/log.h>

@import Appcore;

@interface CMNotificationsDelegate ()

@property(nonatomic, weak) id<UNUserNotificationCenterDelegate> originalDelegate;
@property(nonatomic, weak) CriticalMoments *cm;

@end

@implementation CMNotificationsDelegate

- (id)initWithOriginalDelegate:(id<UNUserNotificationCenterDelegate>)originalDelegate andCm:(CriticalMoments *)cm {
    self = [super init];
    if (self) {
        self.originalDelegate = originalDelegate;
        self.cm = cm;
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
    // Log an event for this notification tap.
    NSString *notificationEventId =
        [NSString stringWithFormat:@"notifications:tapped:%@", response.notification.request.identifier];
    [self.cm sendEvent:notificationEventId];

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
        os_log_error(OS_LOG_DEFAULT,
                     "CriticalMoments: Setup Issue\n\nThe CM notification delegate is not registered. As a result, "
                     "tapping notifications from CM will not trigger the correct action and notifications may fire "
                     "more than once. This is likely because a "
                     "custom UNUserNotificationCenterDelegate was registered after starting CriticalMoments.\n\nTo "
                     "resolve, register your custom UNUserNotificationCenterDelegate before calling "
                     "`CriticalMoments.shared.start()`. Your delegate will still be called for any notification not "
                     "triggered by Critical Moments.\n\nThis warning log is only in debug builds.");
    }
}
#endif

@end
