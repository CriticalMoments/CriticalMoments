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

    UNNotificationContent *content = [CMNotificationHandler buildNotificationContent:notifSchedule.notification];

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

+ (UNNotificationContent *)buildNotificationContent:(DatamodelNotification *)notification {
    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
    content.title = notification.title;
    content.body = notification.body;

    if ([@"default" isEqualToString:notification.sound]) {
        content.sound = [UNNotificationSound defaultSound];
    } else if (notification.sound.length > 0) {
        // Note: invalid names fallback to default. Keep system behaviour.
        content.sound = [UNNotificationSound soundNamed:notification.sound];
    }

    if (@available(iOS 15.0, *)) {
        if ([notification hasRelevanceScore]) {
            content.relevanceScore = [notification getRelevanceScore];
        }
    }

    if (@available(iOS 15.0, *)) {
        if ([@"passive" isEqualToString:notification.interruptionLevel]) {
            content.interruptionLevel = UNNotificationInterruptionLevelPassive;
        } else if ([@"active" isEqualToString:notification.interruptionLevel]) {
            content.interruptionLevel = UNNotificationInterruptionLevelActive;
        } else if ([@"critical" isEqualToString:notification.interruptionLevel]) {
            content.interruptionLevel = UNNotificationInterruptionLevelCritical;
        } else if ([@"timeSensitive" isEqualToString:notification.interruptionLevel]) {
            content.interruptionLevel = UNNotificationInterruptionLevelTimeSensitive;
        }
    }

    return content;
}

+ (NSString *)notificationId:(DatamodelNotification *)notif {
    return [NSString stringWithFormat:@"io.criticalmoments.notifications.%@", notif.id_];
}

@end
