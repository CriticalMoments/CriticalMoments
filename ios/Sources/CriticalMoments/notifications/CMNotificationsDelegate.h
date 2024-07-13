//
//  CMNotificationsDelegate.h
//
//
//  Created by Steve Cosman on 2024-07-06.
//

#import <Foundation/Foundation.h>

#import <UserNotifications/UserNotifications.h>

NS_ASSUME_NONNULL_BEGIN

@interface CMNotificationsDelegate : NSObject <UNUserNotificationCenterDelegate>

- (instancetype)init NS_UNAVAILABLE;
- (instancetype)initWithOriginalDelegate:(id<UNUserNotificationCenterDelegate>)originalDelegate;

#ifdef DEBUG
+ (void)devModeCheckNotificationDelegateSetupCorrectly;
#endif

@end

NS_ASSUME_NONNULL_END
