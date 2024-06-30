//
//  CMNotificationHandler.h
//
//
//  Created by Steve Cosman on 2024-06-30.
//

#import <Foundation/Foundation.h>

@import Appcore;

NS_ASSUME_NONNULL_BEGIN

@interface CMNotificationHandler : NSObject

+ (void)updateNotificationPlan:(AppcoreNotificationPlan *_Nullable)notifPlan;

@end

NS_ASSUME_NONNULL_END
