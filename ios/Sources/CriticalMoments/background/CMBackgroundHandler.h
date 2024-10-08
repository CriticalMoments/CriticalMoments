//
//  CMBackgroundHandler.h
//
//
//  Created by Steve Cosman on 2024-07-11.
//

#import <Foundation/Foundation.h>

#import "../include/CriticalMoments.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMBackgroundHandler : NSObject

- (instancetype)init NS_UNAVAILABLE;
- (instancetype)initWithCm:(CriticalMoments *)cm;
- (void)registerBackgroundTasks;
- (void)scheduleBackgroundTaskAtEpochTime:(int64_t)epochTime;

#ifdef DEBUG
+ (void)devModeCheckBackgroundSetupCorrectly;
#endif

@end

NS_ASSUME_NONNULL_END
