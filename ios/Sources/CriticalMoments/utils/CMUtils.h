//
//  CMUtils.h
//
//
//  Created by Steve Cosman on 2023-05-03.
//

#import <Foundation/Foundation.h>

#define CM_LIB_VERSION_NUMBER_STRING @"0.1.4-beta"

@import UIKit;

NS_ASSUME_NONNULL_BEGIN

@interface CMUtils : NSObject

+ (UIColor *)colorFromHexString:(NSString *)hexString;
+ (UIViewController *)topViewController;
+ (UIWindow *)keyWindow;
+ (NSString *)uiKitLocalizedStringForKey:(NSString *)key;
+ (long long)cmTimestampFromDate:(NSDate *)date;

@end

NS_ASSUME_NONNULL_END
