//
//  CMNotificationHandler.m
//
//
//  Created by Steve Cosman on 2024-06-30.
//

#import "CMNotificationHandler.h"

#import "UserNotifications/UserNotifications.h"

@interface CMNotificationHandler ()

@property(nonatomic, weak) CriticalMoments *cm;

@end

@implementation CMNotificationHandler

- (id)initWithCm:(CriticalMoments *)cm {
    self = [super init];
    if (self) {
        self.cm = cm;
    }
    return self;
}

- (void)updateNotificationPlan:(AppcoreNotificationPlan *_Nullable)plan {
    // Schedule needed notifications
    NSMutableSet<NSString *> *scheduleNotifIds = [[NSMutableSet alloc] init];
    for (int i = 0; i < plan.scheduledNotificationCount; i++) {
        AppcoreScheduledNotification *sn = [plan scheduledNotificationAtIndex:i];
        [self scheduleNotification:sn];
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
              // Notifications in our namespace, which aren't scheduled
              if (![scheduleNotifIds containsObject:request.identifier] &&
                  [request.identifier hasPrefix:DatamodelNotificationUniqueIDPrefix]) {
                  [notifIdsToUnschedule addObject:request.identifier];
              }
          }
          [center removePendingNotificationRequestsWithIdentifiers:notifIdsToUnschedule];
        }];
}

- (void)scheduleNotification:(AppcoreScheduledNotification *)notifSchedule {
    if (!notifSchedule) {
        return;
    }

    UNNotificationContent *content = [self buildNotificationContent:notifSchedule.notification];

    NSDate *date = [NSDate dateWithTimeIntervalSince1970:notifSchedule.scheduledAtEpochMilliseconds / 1000.0];
    UNNotificationTrigger *trigger = [self triggerForDate:date withNotif:notifSchedule.notification];
    if (!trigger) {
        // No error, nil here means already delivered
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

- (UNNotificationTrigger *)triggerForDate:(NSDate *)date withNotif:(DatamodelNotification *)notification {
    NSTimeInterval timeUntilDate = [date timeIntervalSinceNow];
    // Appcore sends all notifications. Including old ones which may be delivered (could be 6 months old). Don't
    // schedule a notification if it's in the past (by more than 0.8s). Those are likely already delivered/scheduled, if
    // not stale.
    if (timeUntilDate < -0.8) {
        return nil;
    }
    if (timeUntilDate <= 1) {
        // Part 1) <= 0s delay not allowed, so check and set a positive value
        // Part 2) Why 1s in the future? In case AppCore sends several updates for same notification rapidly. By
        // scheduling 1s out (from delivery time), and not scheduling if 0.9s in past we avoid duplicate notifications
        // for same timestamp. Could add a cache, but async APIs make that risky. Keep it simple.
        // Part 3) A bit of debounce if app sends multiple events
        timeUntilDate = 1;

        // Log an event for this notification delivery.
        // Note: this isn't guaranteed to be fired for every event, as they can be sent when in the background.
        // **Critical**: this event must be sent for any "now" delivery. Our backend go logic relies on having it as a
        // record that a notification was delivered. Particularly in the case where an "ideal time" allows the
        // notification to be delivered early, we use this event to prevent repeat deliveries.
        NSString *deliveredEventName = [notification deliveredEventName];
        [self.cm sendEvent:deliveredEventName];
    }
    return [UNTimeIntervalNotificationTrigger triggerWithTimeInterval:timeUntilDate repeats:NO];
}

- (UNNotificationContent *)buildNotificationContent:(DatamodelNotification *)notification {
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
