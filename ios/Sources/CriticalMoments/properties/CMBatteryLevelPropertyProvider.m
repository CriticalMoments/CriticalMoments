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

    float batteryLevel = UIDevice.currentDevice.batteryLevel;
    // wired computers and simulators report -1. Let's return 1.0 since their battery is effectively "full".
    if (batteryLevel == -1.0) {
        return 1.0;
    }
    return batteryLevel;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeFloat;
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

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

@end

@implementation CMLowPowerModePropertyProvider

- (BOOL)boolValue {
    NSProcessInfo *processInfo = [[NSProcessInfo alloc] init];
    return processInfo.lowPowerModeEnabled;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

@end
