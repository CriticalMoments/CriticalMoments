//
//  CMViewPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-05-24.
//

#import "CMViewPropertyProvider.h"

#import "../utils/CMUtils.h"

@import UIKit;

@implementation CMDeviceOrientationPropertyProvider

- (NSString *)stringValue {
    UIDeviceOrientation orientation = UIDevice.currentDevice.orientation;
    if (UIDeviceOrientationIsLandscape(orientation)) {
        return @"landscape";
    } else if (UIDeviceOrientationIsPortrait(orientation)) {
        return @"portrait";
    } else if (orientation == UIDeviceOrientationFaceUp) {
        return @"face_up";
    } else if (orientation == UIDeviceOrientationFaceDown) {
        return @"face_down";
    }
    // simulator and first 0.1s after launch are unknown
    return @"unknown";
}

- (long)type {
    return AppcoreLibPropertyProviderTypeString;
}

@end

@implementation CMInterfaceOrientationPropertyProvider

- (NSString *)stringValue {
    // screen reflects the UI
    CGSize screenSize = UIScreen.mainScreen.bounds.size;
    if (screenSize.width > screenSize.height) {
        return @"landscape";
    }
    return @"portrait";
}

- (long)type {
    return AppcoreLibPropertyProviderTypeString;
}

@end

@implementation CMDarkModePropertyProvider

- (BOOL)boolValue {
    if (@available(iOS 12.0, *)) {
        UITraitCollection *tc = UIScreen.mainScreen.traitCollection;
        if (tc.userInterfaceStyle == UIUserInterfaceStyleDark) {
            return YES;
        }
    }
    return NO;
}

- (long)type {
    return AppcoreLibPropertyProviderTypeBool;
}

@end
