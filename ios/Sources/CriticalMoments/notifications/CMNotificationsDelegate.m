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
    // Prefer original delegate behaviour if set
    if (self.originalDelegate) {
        [self.originalDelegate userNotificationCenter:center
                              willPresentNotification:notification
                                withCompletionHandler:completionHandler];
        return;
    }

    // Align to the default OS behaviour if the app doesn't set a delegate: don't display any notifications.
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

    if (self.originalDelegate) {
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

@end
