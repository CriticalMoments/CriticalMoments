//
//  Utils.h
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import <Foundation/Foundation.h>

@import UIKit;

NS_ASSUME_NONNULL_BEGIN

@interface Utils : NSObject

+ (UIWindow *)keyWindow;
+ (UINavigationController *)appNavControl;
+ (void)createTestFileUrls;
+ (BOOL)verifyTestFileUrls;
+ (UIColor *)colorFromHexString:(NSString *)hexString;

@end

NS_ASSUME_NONNULL_END
