//
//  CMUtils.m
//
//
//  Created by Steve Cosman on 2023-05-03.
//

#import "CMUtils.h"
#import "../CriticalMoments_private.h"
#import "../include/CriticalMoments.h"

@import Appcore;

@implementation CMUtils

/// Parse hex codes in format #ffffff to UIColor
+ (UIColor *)colorFromHexString:(NSString *)hexString {
    if (hexString.length != 7) {
        return nil;
    }

    unsigned int parsed = 0;
    NSScanner *scanner = [NSScanner scannerWithString:hexString];

    if ([hexString hasPrefix:@"#"]) {
        [scanner setScanLocation:1];
    } else {
        return nil;
    }
    bool scannedHex = [scanner scanHexInt:&parsed];
    if (!scannedHex || ![scanner isAtEnd]) {
        return nil;
    }

    CGFloat red = ((parsed & 0xff0000) >> 16) / 255.0;
    CGFloat green = ((parsed & 0x00ff00) >> 8) / 255.0;
    CGFloat blue = (parsed & 0x0000ff) / 255.0;

    return [[UIColor alloc] initWithRed:red green:green blue:blue alpha:1.0];
}

+ (UIWindow *)keyWindow {
    UIWindow *keyWindow = [[[UIApplication sharedApplication] windows] firstObject];
    for (UIWindow *w in [[UIApplication sharedApplication] windows]) {
        if (w.isKeyWindow) {
            keyWindow = w;
            break;
        }
    }
    return keyWindow;
}

+ (UIViewController *)topViewController {
    UIViewController *vc = [CMUtils keyWindow].rootViewController;

    // Find top VC, unless it's in the process of being dismissed
    for (UIViewController *nextPresented = vc.presentedViewController; nextPresented && !nextPresented.beingDismissed;
         nextPresented = vc.presentedViewController) {
        vc = nextPresented;
    }

    return vc;
}

+ (NSString *)uiKitLocalizedStringForKey:(NSString *)key {
    NSBundle *uikitBundle = [NSBundle bundleForClass:[UIButton class]];
    if (!uikitBundle) {
        return key;
    }
    return [uikitBundle localizedStringForKey:key value:key table:nil];
}

+ (long long)cmTimestampFromDate:(NSDate *)date {
    NSTimeInterval unixTime = [date timeIntervalSince1970];
    return unixTime * 1000;
}

+ (bool)isiPad {
    return UIDevice.currentDevice.userInterfaceIdiom == UIUserInterfaceIdiomPad;
}

+ (int64_t)dateToGoTime:(NSDate *)value {
    if (!value) {
        return AppcoreLibPropertyProviderNilIntValue;
    } else {
        int64_t epochMilliseconds = [@(floor([value timeIntervalSince1970] * 1000)) longLongValue];
        return epochMilliseconds;
    }
}

+ (NSDictionary *)fetchCmApiSyncronous:(NSString *)urlString error:(NSError **)error {
    NSURL *url = [NSURL URLWithString:urlString];
    NSMutableURLRequest *req = [[NSMutableURLRequest alloc] initWithURL:url];
    // TODO P2 - don't use shared instance here, move API into CM instance for testing
    NSString *apiKey = [CriticalMoments.sharedInstance getApiKey];
    [req setValue:apiKey forHTTPHeaderField:@"X-CM-Api-Key"];

    __block NSDictionary *jsonDict;
    __block NSError *statusErr;

    dispatch_semaphore_t approxSem = dispatch_semaphore_create(0);
    [[[NSURLSession sharedSession]
        dataTaskWithRequest:req
          completionHandler:^(NSData *data, NSURLResponse *response, NSError *reqErr) {
            NSHTTPURLResponse *httpResp;
            if ([response isKindOfClass:[NSHTTPURLResponse class]]) {
                httpResp = (NSHTTPURLResponse *)response;
            }
            if (httpResp.statusCode != 200) {
                statusErr = [[NSError alloc] initWithDomain:@"io.criticalmoments"
                                                       code:21345123532
                                                   userInfo:@{@"message" : @"API Status Code != 200"}];
                // continue to parse JSON for error message
            }
            if (!reqErr && data.length > 0 && httpResp) {
                NSError *error = nil;
                id jsonObj = [NSJSONSerialization JSONObjectWithData:data options:0 error:&error];
                if (!error && [jsonObj isKindOfClass:[NSDictionary class]]) {
                    jsonDict = (NSDictionary *)jsonObj;
                }
            }
            dispatch_semaphore_signal(approxSem);
          }] resume];
    dispatch_semaphore_wait(approxSem, dispatch_time(DISPATCH_TIME_NOW, 10.0 * NSEC_PER_SEC));

    // copy point in time pointer because of async semaphore above
    NSDictionary *returnDict = jsonDict;
    if (returnDict && returnDict[@"errorMessage"] != nil) {
        *error = [[NSError alloc] initWithDomain:@"io.criticalmoments"
                                            code:3784523948
                                        userInfo:@{@"message" : returnDict[@"errorMessage"]}];
        return nil;
    }
    if (statusErr) {
        *error = statusErr;
        return nil;
    }
    if (!returnDict) {
        *error = [[NSError alloc] initWithDomain:@"io.criticalmoments"
                                            code:2348790234
                                        userInfo:@{@"message" : @"API timeout or invalid content"}];
        return nil;
    }

    return returnDict;
}

#define MAX_BEARING (2.0 * M_PI)

+ (CLLocation *)noiseLocation:(CLLocation *)location maxNoise:(int)distanceInMeters {
    // random bearing and distance offset
    float distanceMeters = ((float)arc4random() / UINT32_MAX) * distanceInMeters;
    float radianBearing = ((float)arc4random() / UINT32_MAX) * MAX_BEARING;

    const double distRadians = distanceMeters / (6372797.6); // earth radius in meters

    float lat1 = location.coordinate.latitude * M_PI / 180.0;
    float lon1 = location.coordinate.longitude * M_PI / 180.0;

    float lat2 = asin(sin(lat1) * cos(distRadians) + cos(lat1) * sin(distRadians) * cos(radianBearing));
    float lon2 =
        lon1 + atan2(sin(radianBearing) * sin(distRadians) * cos(lat1), cos(distRadians) - sin(lat1) * sin(lat2));

    CLLocationDegrees finalLat = lat2 * 180.0 / M_PI;
    if (finalLat > 90.0) {
        finalLat = 90.0;
    } else if (finalLat < -90.0) {
        finalLat = -90.0;
    }
    CLLocationDegrees finalLong = lon2 * 180.0 / M_PI;
    if (finalLong > 180.0) {
        finalLong = 180.0;
    } else if (finalLong < -180.0) {
        finalLong = -180.0;
    }
    return [[CLLocation alloc] initWithLatitude:finalLat longitude:finalLong];
}

@end
