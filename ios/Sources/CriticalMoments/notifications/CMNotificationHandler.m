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
    NSMutableSet<NSString *> *scheduleNotifIds = [[NSMutableSet alloc] init];
    for (int i = 0; i < plan.scheduledNotificationCount; i++) {
        AppcoreScheduledNotification *sn = [plan scheduledNotificationAtIndex:i];
        [CMNotificationHandler scheduleNotification:sn];
        [scheduleNotifIds addObject:[sn.notification uniqueID]];
    }

    // Unschedule all notifications that are not in scheduled list
    // Note: can't simply use plan.unscheduledNotification, as we also want to delete any from prior configs. Use our
    // namespace instead.
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center
        getPendingNotificationRequestsWithCompletionHandler:^(NSArray<UNNotificationRequest *> *_Nullable requests) {
          NSMutableArray<NSString *> *notifIdsToUnschedule = [[NSMutableArray alloc] init];
          for (UNNotificationRequest *request in requests) {
              // Notificaitons in our namespace, which aren't scheduled
              if (![scheduleNotifIds containsObject:request.identifier] &&
                  [request.identifier hasPrefix:DatamodelNotificationUniqueIDPrefix]) {
                  [notifIdsToUnschedule addObject:request.identifier];
              }
          }
          [center removePendingNotificationRequestsWithIdentifiers:notifIdsToUnschedule];
        }];
}

+ (void)scheduleNotification:(AppcoreScheduledNotification *)notifSchedule {
    if (!notifSchedule) {
        return;
    }

    UNNotificationContent *content = [CMNotificationHandler buildNotificationContent:notifSchedule.notification];

    NSDate *date = [NSDate dateWithTimeIntervalSince1970:notifSchedule.scheduledAtEpochMilliseconds / 1000.0];
    UNNotificationTrigger *trigger = [CMNotificationHandler triggerForDate:date];
    if (!trigger) {
        return;
    }

    NSString *notifId = [notifSchedule.notification uniqueID];

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

+ (UNNotificationTrigger *)triggerForDate:(NSDate *)date {
    NSTimeInterval timeUntilDate = [date timeIntervalSinceNow];
    // TODO_P0: the -0.5 check is checking that this was just scheduled (not old and likely already delivered one). Temp
    // workaround. We should cache delivered time (or scheduled time) by ID, and not re-schedule if we already
    // delivered/scheduled it.
    if (timeUntilDate < -0.5) {
        return nil;
    }
    if (timeUntilDate <= 0.5) {
        // <= 0s delay not allowed, so use small value.
        timeUntilDate = 0.001;
    }
    return [UNTimeIntervalNotificationTrigger triggerWithTimeInterval:timeUntilDate repeats:NO];
}

+ (UNNotificationContent *)buildNotificationContent:(DatamodelNotification *)notification {
    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];

    content.title = notification.title;
    content.body = notification.body;

    // Zero is valid to remove badge. Any negative value should be nil.
    if (notification.badgeCount >= 0) {
        content.badge = [NSNumber numberWithLong:notification.badgeCount];
    }

    if (notification.launchImageName.length > 0) {
        content.launchImageName = notification.launchImageName;
    }

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

@end
