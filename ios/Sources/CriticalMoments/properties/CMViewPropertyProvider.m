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

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
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

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
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

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

@end

@implementation CMAppStatePropertyProvider

+ (UIApplicationState)getApplicationStateBlocking {
    // [UIApplication applicationState] must be used from main thread only
    UIApplicationState __block state = -999999;
    dispatch_semaphore_t stateSem = dispatch_semaphore_create(0);
    dispatch_async(dispatch_get_main_queue(), ^{
      state = UIApplication.sharedApplication.applicationState;
      dispatch_semaphore_signal(stateSem);
    });
    dispatch_semaphore_wait(stateSem, dispatch_time(DISPATCH_TIME_NOW, 0.5 * NSEC_PER_SEC));
    return state;
}

- (NSString *)stringValue {
    UIApplicationState state = [CMAppStatePropertyProvider getApplicationStateBlocking];
    switch (state) {
    case UIApplicationStateActive:
        return @"active";
    case UIApplicationStateInactive:
        return @"inactive";
    case UIApplicationStateBackground:
        return @"background";
    }

    return @"unknown";
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

@end

@implementation CMForegroundProvider

- (BOOL)boolValue {
    UIApplicationState state = [CMAppStatePropertyProvider getApplicationStateBlocking];
    return state != UIApplicationStateBackground;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

@end

@implementation CMBrightnessProvider

- (double)floatValue {
    return UIScreen.mainScreen.brightness;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeFloat;
}

@end
