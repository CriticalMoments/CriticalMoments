//
//  CMUtils.h
//
//
//  Created by Steve Cosman on 2023-05-03.
//

#import <Foundation/Foundation.h>

#define CM_LIB_VERSION_NUMBER_STRING @"0.9.5"

@import UIKit;
@import CoreLocation;

NS_ASSUME_NONNULL_BEGIN

@interface CMUtils : NSObject

+ (UIColor *)colorFromHexString:(NSString *)hexString;
+ (UIViewController *)topViewController;
+ (UIWindow *)keyWindow;
+ (NSString *)uiKitLocalizedStringForKey:(NSString *)key;
+ (long long)cmTimestampFromDate:(NSDate *)date;
+ (bool)isiPad;
+ (int64_t)dateToGoTime:(NSDate *)value;
+ (NSDictionary *)fetchCmApiSyncronous:(NSString *)urlString error:(NSError **)error;
+ (CLLocation *)noiseLocation:(CLLocation *)loc maxNoise:(int)distanceInMeters;

@end

NS_ASSUME_NONNULL_END
