//
//  CMUtils.h
//
//
//  Created by Steve Cosman on 2023-05-03.
//

#import <Foundation/Foundation.h>

#define CM_LIB_VERSION_NUMBER_STRING @"0.2.3-beta"

@import UIKit;

NS_ASSUME_NONNULL_BEGIN

@interface CMUtils : NSObject

+ (UIColor *)colorFromHexString:(NSString *)hexString;
+ (UIViewController *)topViewController;
+ (UIWindow *)keyWindow;
+ (NSString *)uiKitLocalizedStringForKey:(NSString *)key;
+ (long long)cmTimestampFromDate:(NSDate *)date;
+ (bool)isiPad;

@end

NS_ASSUME_NONNULL_END
