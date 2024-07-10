//
//  NotificationsDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2024-07-06.
//

#import "NotificationsDemoScreen.h"

#import "UserNotifications/UserNotifications.h"
#import <CriticalMoments.h>
#import <UIKit/UIKit.h>

@implementation NotificationsDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Notification & Badges";
        self.infoText =
            @"Deliver notifiations to users when your app isn't open.\n\nSet the badge (count) on your app "
            @"icon.\n\nCritical Moments can deliver "
            @"notifications without servers. It can also optimize the ideal delivery time to increase engagement.";
        // self.buttonLink = @"https://docs.criticalmoments.io/actions/alerts";

        [self buildSections];
    }
    return self;
}

- (void)didAppear:(UIViewController *)vc {
    [self requestPermission:vc];
}

- (void)buildSections {
    // Future usage demos:
    // Remind later: when on WiFi and using device. Note: you won't recieve this notificaiton right away, it make
    // take a few hours. But it is likely to find you when you're at home/work, and have a moment to use your phone.

    CMDemoAction *comeBackDemo = [[CMDemoAction alloc] init];
    comeBackDemo.title = @"Complete Onboarding CTA";
    comeBackDemo.subtitle =
        @"Reduce new-user churn by reminding users who abandon setup. Show a notification a few minutes after they "
        @"leave the app (for this demo we use a 5s delay).\n\nTap this, then dismiss the app to see a CTA notification "
        @"5 seconds after the app moves to the background.";
    comeBackDemo.actionCMEventName = @"enableComeBack";

    [self addSection:@"Use Case Demos" withActions:@[ comeBackDemo ]];

    CMDemoAction *eventNotification = [[CMDemoAction alloc] init];
    eventNotification.title = @"Basic Notifications";
    eventNotification.subtitle = @"Display a notification, triggered immediatly after an in app event.";
    eventNotification.actionCMEventName = @"demo_notification_1";

    CMDemoAction *delayNotification = [[CMDemoAction alloc] init];
    delayNotification.title = @"Delayed Notifications";
    delayNotification.subtitle = @"Display a notification with a 5s delay.";
    delayNotification.actionCMEventName = @"demo_notification_2";

    CMDemoAction *badgeNotif = [[CMDemoAction alloc] init];
    badgeNotif.title = @"Set App Icon Badge";
    badgeNotif.subtitle = @"Display a badge on your app icon. Badges are numbers that indicate content is available in "
                          @"app. Open your homescreen after tapping to see the effect.";
    badgeNotif.actionCMEventName = @"demo_notification_3";

    CMDemoAction *clearBadgeNotif = [[CMDemoAction alloc] init];
    clearBadgeNotif.title = @"Clear App Icon Badge";
    clearBadgeNotif.subtitle =
        @"After setting with demo above, tap here to remove the badge from your app icon. Open your homescreen "
        @"after tapping to see the effect.";
    clearBadgeNotif.actionCMEventName = @"demo_notification_4";

    CMDemoAction *criticalNotif = [[CMDemoAction alloc] init];
    criticalNotif.title = @"Critical Notification";
    criticalNotif.subtitle =
        @"Show a notification with a 'critical' interruption level, which will receive priority on the lock screen.";
    criticalNotif.actionCMEventName = @"demo_notification_5";

    [self addSection:@"Notification Examples" withActions:@[ eventNotification, delayNotification, criticalNotif ]];

    [self addSection:@"Badge Examples" withActions:@[ badgeNotif, clearBadgeNotif ]];
}

- (void)requestPermission:(UIViewController *)vc {
    [CriticalMoments.shared
        requestNotificationPermissionWithCompletionHandler:^(BOOL granted, NSError *_Nullable error) {
          if (!granted || error) {
              [self showPermissionsIssue:vc];
          }
        }];
}

- (void)showPermissionsIssue:(UIViewController *)vc {
    dispatch_async(dispatch_get_main_queue(), ^{
      // Always go back since demos on this screen don't work without permissions
      [vc.navigationController popViewControllerAnimated:YES];

      UIAlertController *alert =
          [UIAlertController alertControllerWithTitle:@"Notification Permission Required"
                                              message:@"This demo needs permission to deliver notifications. Please "
                                                      @"allow permissions in Settings."
                                       preferredStyle:UIAlertControllerStyleAlert];

      UIAlertAction *cancelAction = [UIAlertAction actionWithTitle:@"Cancel"
                                                             style:UIAlertActionStyleCancel
                                                           handler:^(UIAlertAction *action){
                                                           }];
      [alert addAction:cancelAction];
      UIAlertAction *defaultAction =
          [UIAlertAction actionWithTitle:@"Settings"
                                   style:UIAlertActionStyleDefault
                                 handler:^(UIAlertAction *action) {
                                   NSURL *url = [[NSURL alloc] initWithString:UIApplicationOpenSettingsURLString];
                                   [[UIApplication sharedApplication] openURL:url options:@{} completionHandler:nil];
                                 }];
      [alert addAction:defaultAction];
      alert.preferredAction = defaultAction;
      [vc presentViewController:alert animated:YES completion:nil];
    });
}

@end