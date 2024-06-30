//
//  CMNotificationHandler.m
//
//
//  Created by Steve Cosman on 2024-06-30.
//

#import "CMNotificationHandler.h"

#import "UserNotifications/UserNotifications.h"

@implementation CMNotificationHandler

+ (void)updateNotificationPlan:(AppcoreNotificationPlan *_Nullable)plan {
    // Schedule needed notifications
    for (int i = 0; i < plan.scheduledNotificationCount; i++) {
        AppcoreScheduledNotification *notif = [plan scheduledNotificationAtIndex:i];
        [CMNotificationHandler scheduleNotification:notif];
    }

    // Unschedule all notifications that no longer apply
    if (plan.unscheduledNotificationCount > 0) {
        NSMutableArray<NSString *> *unscheduleNotifIds = [[NSMutableArray alloc] init];
        for (int i = 0; i < plan.unscheduledNotificationCount; i++) {
            DatamodelNotification *notif = [plan unscheduledNotificationAtIndex:i];
            if (!notif) {
                continue;
            }
            NSString *notifId = [CMNotificationHandler notificationId:notif];
            [unscheduleNotifIds addObject:notifId];
        }
        UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
        [center removePendingNotificationRequestsWithIdentifiers:unscheduleNotifIds];
    }
}

+ (void)scheduleNotification:(AppcoreScheduledNotification *)notifSchedule {
    if (!notifSchedule) {
        return;
    }

    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
    content.title = notifSchedule.notification.title;
    // TODO_P0 check this can be empty, we allow empty in config
    content.body = notifSchedule.notification.body;
    // TODO_P0 -- more options, at least silent/default.
    // content.sound = [UNNotificationSound defaultSound];

    // TODO: would really rather use exact time. This could be messed up by timezones? And too complicated. Look at
    // android API for how I want to implement dow and tod filters (go or here). Here might respect timezone!
    NSDate *date = [NSDate dateWithTimeIntervalSince1970:notifSchedule.scheduledAtEpochMilliseconds / 1000.0];
    NSCalendarUnit allUnits = NSCalendarUnitYear | NSCalendarUnitMonth | NSCalendarUnitDay | NSCalendarUnitHour |
                              NSCalendarUnitMinute | NSCalendarUnitSecond | NSCalendarUnitTimeZone;
    NSDateComponents *dateComponents = [NSCalendar.currentCalendar components:allUnits fromDate:date];

    UNCalendarNotificationTrigger *trigger =
        [UNCalendarNotificationTrigger triggerWithDateMatchingComponents:dateComponents repeats:NO];

    NSString *notifId = [CMNotificationHandler notificationId:notifSchedule.notification];

    UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:notifId
                                                                          content:content
                                                                          trigger:trigger];

    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center addNotificationRequest:request
             withCompletionHandler:^(NSError *_Nullable error) {
               if (error) {
                   NSLog(@"CriticalMoments: Error scheduling notification: %@", error);
               }
             }];
}

+ (NSString *)notificationId:(DatamodelNotification *)notif {
    return [NSString stringWithFormat:@"io.criticalmoments.notifications.%@", notif.id_];
}

@end
