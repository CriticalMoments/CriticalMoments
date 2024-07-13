//
//  CMBackgroundHandler.m
//
//
//  Created by Steve Cosman on 2024-07-11.
//

#import "CMBackgroundHandler.h"

#import <BackgroundTasks/BackgroundTasks.h>

// TODO_P0 remove import
#import "UserNotifications/UserNotifications.h"

#define bgFetchTaskId @"io.criticalmoments.bg_fetch"
#define bgProcessingTaskId @"io.criticalmoments.bg_process"
#define allBackgroundIds @[ bgFetchTaskId, bgProcessingTaskId ]

@implementation CMBackgroundHandler

+ (void)registerBackgroundTasks {
    if (@available(iOS 13.0, *)) {
        // Simulators do not support background work
#ifndef TARGET_IPHONE_SIMULATOR
        for (NSString *taskId in allBackgroundIds) {
            BOOL registered =
                [BGTaskScheduler.sharedScheduler registerForTaskWithIdentifier:taskId
                                                                    usingQueue:nil
                                                                 launchHandler:^(__kindof BGTask *_Nonnull task) {
                                                                   [CMBackgroundHandler runBackgroundWorker:task];
                                                                 }];
            if (!registered) {
                [CMBackgroundHandler logSetupError:taskId];
            }
        }
#endif
    }
}

+ (void)scheduleBackgroundTask {
    // Simulators do not support background work APIs
#ifdef TARGET_IPHONE_SIMULATOR
    return;
#endif

    if (@available(iOS 13.0, *)) {
        BGAppRefreshTaskRequest *fetchRequest = [[BGAppRefreshTaskRequest alloc] initWithIdentifier:bgFetchTaskId];
        // At least 15 mins from now
        fetchRequest.earliestBeginDate = [NSDate dateWithTimeIntervalSinceNow:15 * 60];

        NSError *error;
        BOOL success = [BGTaskScheduler.sharedScheduler submitTaskRequest:fetchRequest error:&error];
        if (!success || error) {
            [CMBackgroundHandler logSetupError:bgFetchTaskId];
        }

        error = nil;
        BGProcessingTaskRequest *processingRequest =
            [[BGProcessingTaskRequest alloc] initWithIdentifier:bgProcessingTaskId];
        success = [BGTaskScheduler.sharedScheduler submitTaskRequest:processingRequest error:&error];
        if (!success || error) {
            [CMBackgroundHandler logSetupError:bgProcessingTaskId];
        }
    }

    for (NSString *taskId in allBackgroundIds) {
        NSURL *logPath = [CMBackgroundHandler logPath:true withTaskId:taskId];
        NSString *logContents = [NSString stringWithContentsOfURL:logPath encoding:NSUTF8StringEncoding error:nil];
        NSLog(@"BG Debug Log [%@]:\n%@\n\n", taskId, logContents);

        logPath = [CMBackgroundHandler logPath:false withTaskId:taskId];
        logContents = [NSString stringWithContentsOfURL:logPath encoding:NSUTF8StringEncoding error:nil];
        NSLog(@"BG Release Log [%@]:\n%@\n\n", taskId, logContents);
    }
}

+ (void)runBackgroundWorker:(BGTask *)task API_AVAILABLE(ios(13.0)) {
    // Schedule next refresh
    [CMBackgroundHandler scheduleBackgroundTask];

    [CMBackgroundHandler logRunTimestamp:task.identifier];
    [CMBackgroundHandler scheduleNotificationNow:task.identifier];
    NSLog(@"CMBackground: worker ran - %@", task.identifier);

    [task setTaskCompletedWithSuccess:YES];
}

// TODO_P0 remove this
+ (void)scheduleNotificationNow:(NSString *)taskId {
    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];

    if ([bgFetchTaskId isEqualToString:taskId]) {
        content.title = @"Background Fetch";
    } else {
        content.title = @"Background Processing";
    }
    NSString *dateString = [NSDateFormatter localizedStringFromDate:[NSDate date]
                                                          dateStyle:NSDateFormatterShortStyle
                                                          timeStyle:NSDateFormatterFullStyle];

    content.body = dateString;

    UNNotificationTrigger *trigger = [UNTimeIntervalNotificationTrigger triggerWithTimeInterval:1.0 repeats:NO];

    UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:@"notif_bg"
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

// TODO_P0 remove this
+ (NSURL *)logPath:(BOOL)debug withTaskId:(NSString *)taskId {
    NSURL *appSupportDir = [[NSFileManager.defaultManager URLsForDirectory:NSApplicationSupportDirectory
                                                                 inDomains:NSUserDomainMask] lastObject];

    NSURL *criticalMomentsDataDir = [appSupportDir URLByAppendingPathComponent:@"critical_moments_test_data"];
    NSError *error;
    BOOL s = [NSFileManager.defaultManager createDirectoryAtURL:criticalMomentsDataDir
                                    withIntermediateDirectories:YES
                                                     attributes:nil
                                                          error:&error];
    if (!s || error) {
        NSLog(@"error: %@", error);
    }
    NSString *filename = [NSString stringWithFormat:@"%@.log", taskId];
    if (debug) {
        filename = [NSString stringWithFormat:@"%@_debug.log", taskId];
    }
    NSURL *bgLogFile = [criticalMomentsDataDir URLByAppendingPathComponent:filename];
    return bgLogFile;
}

// TODO_P0 remove this
+ (void)logRunTimestamp:(NSString *)taskId {
    NSString *dateString = [NSDateFormatter localizedStringFromDate:[NSDate date]
                                                          dateStyle:NSDateFormatterShortStyle
                                                          timeStyle:NSDateFormatterFullStyle];

#ifdef DEBUG
    NSURL *logPath = [CMBackgroundHandler logPath:true withTaskId:taskId];
#else
    NSURL *logPath = [CMBackgroundHandler logPath:false withTaskId:taskId];
#endif
    NSString *logContents = [NSString stringWithContentsOfURL:logPath encoding:NSUTF8StringEncoding error:nil];
    NSString *newContent = dateString;
    if (logContents) {
        newContent = [NSString stringWithFormat:@"%@\n%@", logContents, dateString];
    }
    NSError *error;
    BOOL s = [newContent writeToURL:logPath atomically:YES encoding:NSUTF8StringEncoding error:&error];
    if (!s || error) {
        NSLog(@"error: %@", error);
    }
}

+ (void)logSetupError:(NSString *)taskId {
    NSLog(@"CriticalMoments: failed to register background worker [%@]. Please ensure you follow all the steps in our "
          @"quick "
          @"start guide. https://docs.criticalmoments.io/quick-start",
          taskId);
}

#ifdef DEBUG
// Check everything is setup correctly, and log a warning if not.
// Only compiled in debug mode, won't run on release builds.
+ (void)devModeCheckBackgroundSetupCorrectly {
    // Check our 2 IDs are included in the app's Info.plist
    NSArray *permittedIdentifiers =
        [[NSBundle mainBundle] objectForInfoDictionaryKey:@"BGTaskSchedulerPermittedIdentifiers"];
    for (NSString *requiredTaskId in allBackgroundIds) {
        if (![permittedIdentifiers containsObject:requiredTaskId]) {
            NSLog(@"CriticalMoments Debug Developer Warning: The CM background task IDs must be registered in your "
                  @"Info.plist for some Critical Moments features to function, such as notifications. See our quick "
                  @"start guide for details: https://docs.criticalmoments.io\nThis warning is only on debug builds, "
                  @"and is not included in release builds.");
            break;
        }
    }
}
#endif

@end
