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
}

+ (void)runBackgroundWorker:(BGTask *)task API_AVAILABLE(ios(13.0)) {
    // Schedule next refresh
    [CMBackgroundHandler scheduleBackgroundTask];

    NSLog(@"Background worker ran");
    [task setTaskCompletedWithSuccess:YES];
}

+ (void)logSetupError {
    NSLog(@"CriticalMoments: failed to register background worker. Please ensure you follow all the steps in our quick "
          @"start guide. https://docs.criticalmoments.io/quick-start");
}

@end
