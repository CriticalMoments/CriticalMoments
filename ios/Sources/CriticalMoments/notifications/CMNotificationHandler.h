//
//  CMNotificationHandler.h
//
//
//  Created by Steve Cosman on 2024-06-30.
//

#import <Foundation/Foundation.h>

#import "../include/CriticalMoments.h"

@import Appcore;

NS_ASSUME_NONNULL_BEGIN

@interface CMNotificationHandler : NSObject

- (instancetype)init NS_UNAVAILABLE;
- (instancetype)initWithCm:(CriticalMoments *)cm;

- (void)updateNotificationPlan:(AppcoreNotificationPlan *_Nullable)notifPlan;

@end

NS_ASSUME_NONNULL_END
