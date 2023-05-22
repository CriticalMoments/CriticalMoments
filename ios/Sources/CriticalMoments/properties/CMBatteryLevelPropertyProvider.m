//
//  CMBatteryLevelPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-05-22.
//

#import "CMBatteryLevelPropertyProvider.h"

@import UIKit;

@implementation CMBatteryLevelPropertyProvider

- (double)floatValue {
    if (!UIDevice.currentDevice.batteryMonitoringEnabled) {
        UIDevice.currentDevice.batteryMonitoringEnabled = YES;
    }

    return UIDevice.currentDevice.batteryLevel;
}

- (long)type {
    return AppcoreLibPropertyProviderTypeFloat;
}

@end

@implementation CMBatteryStatePropertyProvider

- (NSString *)stringValue {
    if (!UIDevice.currentDevice.batteryMonitoringEnabled) {
        UIDevice.currentDevice.batteryMonitoringEnabled = YES;
    }

    switch (UIDevice.currentDevice.batteryState) {
    case UIDeviceBatteryStateCharging:
        return @"charging";
    case UIDeviceBatteryStateFull:
        return @"full";
    case UIDeviceBatteryStateUnplugged:
        return @"unplugged";
    case UIDeviceBatteryStateUnknown:
        break;
    }
    return @"unknown";
}

- (long)type {
    return AppcoreLibPropertyProviderTypeString;
}

@end

@implementation CMLowPowerModePropertyProvider

- (BOOL)boolValue {
    NSProcessInfo *processInfo = [[NSProcessInfo alloc] init];
    return processInfo.lowPowerModeEnabled;
}

- (long)type {
    return AppcoreLibPropertyProviderTypeBool;
}

@end
