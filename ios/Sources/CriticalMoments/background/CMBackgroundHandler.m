//
//  CMBackgroundHandler.m
//
//
//  Created by Steve Cosman on 2024-07-11.
//

#import "CMBackgroundHandler.h"

#import <BackgroundTasks/BackgroundTasks.h>

#define bgTaskId @"io.criticalmoments"

@implementation CMBackgroundHandler

+ (void)registerBackgroundTasks {
    if (@available(iOS 13.0, *)) {
        BOOL registered =
            [BGTaskScheduler.sharedScheduler registerForTaskWithIdentifier:bgTaskId
                                                                usingQueue:nil
                                                             launchHandler:^(__kindof BGTask *_Nonnull task) {
                                                               [CMBackgroundHandler runBackgroundWorker:task];
                                                             }];
        if (!registered) {
            [CMBackgroundHandler logSetupError];
        }
    }
}

+ (void)scheduleBackgroundTask {
    if (@available(iOS 13.0, *)) {
        BGAppRefreshTaskRequest *request = [[BGAppRefreshTaskRequest alloc] initWithIdentifier:bgTaskId];
        // At least 15 mins from now
        request.earliestBeginDate = [NSDate dateWithTimeIntervalSinceNow:15 * 60];

        NSError *error;
        BOOL success = [BGTaskScheduler.sharedScheduler submitTaskRequest:request error:&error];
        if (!success || error) {
            [CMBackgroundHandler logSetupError];
        }
    }

    NSURL *logPath = [CMBackgroundHandler logPath:true];
    NSString *logContents = [NSString stringWithContentsOfURL:logPath encoding:NSUTF8StringEncoding error:nil];
    NSLog(@"BG Debug Log:\n%@\n\n", logContents);

    logPath = [CMBackgroundHandler logPath:false];
    logContents = [NSString stringWithContentsOfURL:logPath encoding:NSUTF8StringEncoding error:nil];
    NSLog(@"BG Release Log:\n%@\n\n", logContents);
}

+ (void)runBackgroundWorker:(BGTask *)task API_AVAILABLE(ios(13.0)) {
    // Schedule next refresh
    [CMBackgroundHandler scheduleBackgroundTask];

    [CMBackgroundHandler logRunTimestamp];
    NSLog(@"Background worker ran");
    [task setTaskCompletedWithSuccess:YES];
}

+ (NSURL *)logPath:(BOOL)debug {
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
    NSURL *bgLogFile = [criticalMomentsDataDir URLByAppendingPathComponent:@"bg_log.log"];
    if (debug) {
        bgLogFile = [criticalMomentsDataDir URLByAppendingPathComponent:@"bg_log_debug.log"];
    }
    return bgLogFile;
}

+ (void)logRunTimestamp {
    NSString *dateString = [NSDateFormatter localizedStringFromDate:[NSDate date]
                                                          dateStyle:NSDateFormatterShortStyle
                                                          timeStyle:NSDateFormatterFullStyle];

#ifdef DEBUG
    NSURL *logPath = [CMBackgroundHandler logPath:true];
#else
    NSURL *logPath = [CMBackgroundHandler logPath:false];
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

+ (void)logSetupError {
    NSLog(@"CriticalMoments: failed to register background worker. Please ensure you follow all the steps in our quick "
          @"start guide. https://docs.criticalmoments.io/quick-start");
}

@end
