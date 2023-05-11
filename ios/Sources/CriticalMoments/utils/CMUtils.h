//
//  CMUtils.h
//
//
//  Created by Steve Cosman on 2023-05-03.
//

#import <Foundation/Foundation.h>

@import UIKit;

NS_ASSUME_NONNULL_BEGIN

@interface CMUtils : NSObject

+ (UIColor *)colorFromHexString:(NSString *)hexString;
+ (UIWindow *)keyWindow;
+ (NSString *)uiKitLocalizedStringForKey:(NSString *)key;

@end

NS_ASSUME_NONNULL_END
