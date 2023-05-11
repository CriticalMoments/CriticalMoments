//
//  CMUtils.m
//
//
//  Created by Steve Cosman on 2023-05-03.
//

#import "CMUtils.h"

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
    UIWindow *keyWindow =
        [[[UIApplication sharedApplication] windows] firstObject];
    for (UIWindow *w in [[UIApplication sharedApplication] windows]) {
        if (w.isKeyWindow) {
            keyWindow = w;
            break;
        }
    }
    return keyWindow;
}

@end
