//
//  CMDefaultProperties.m
//
//
//  Created by Steve Cosman on 2023-05-20.
//

#import "CMDefaultProperties.h"

#import "CMBatteryLevelPropertyProvider.h"

#import <sys/utsname.h>

@import UIKit;

@implementation CMDefaultProperties

+ (void)registerDefaultPropertiesToAppcore {
    AppcoreAppcore *ac = AppcoreSharedAppcore();

    // This API returns different values on older iPads. Make these
    // consistent (documented in docs)
    NSString *systemOsName = UIDevice.currentDevice.systemName;
    if ([@"iOS" isEqualToString:systemOsName] &&
        UIDevice.currentDevice.userInterfaceIdiom == UIUserInterfaceIdiomPad) {
        systemOsName = @"iPadOS";
    }
    [ac registerStaticStringProperty:@"platform" value:systemOsName];

    [CMDefaultProperties setVersionString:UIDevice.currentDevice.systemVersion
                                forPrefix:@"os"];

    // Locale
    NSLocale *locale = [NSLocale currentLocale];
    [ac registerStaticStringProperty:@"locale_language_code"
                               value:locale.languageCode];
    [ac registerStaticStringProperty:@"locale_country_code"
                               value:locale.countryCode];
    [ac registerStaticStringProperty:@"locale_currency_code"
                               value:locale.currencyCode];

    // Bundle ID
    [ac registerStaticStringProperty:@"app_id"
                               value:NSBundle.mainBundle.bundleIdentifier];

    // App Version
    NSString *appVersion = [NSBundle.mainBundle
        objectForInfoDictionaryKey:@"CFBundleShortVersionString"];
    [CMDefaultProperties setVersionString:appVersion forPrefix:@"app"];

    // Screen size / scale
    CGSize screenSize = UIScreen.mainScreen.bounds.size;
    [ac registerStaticIntProperty:@"screen_width_points"
                            value:MIN(screenSize.width, screenSize.height)];
    [ac registerStaticIntProperty:@"screen_height_points"
                            value:MAX(screenSize.width, screenSize.height)];
    CGFloat screenWidthPixels = screenSize.width * UIScreen.mainScreen.scale;
    CGFloat screenHeightPixels = screenSize.height * UIScreen.mainScreen.scale;
    [ac registerStaticIntProperty:@"screen_width_pixels"
                            value:MIN(screenHeightPixels, screenWidthPixels)];
    [ac registerStaticIntProperty:@"screen_height_pixels"
                            value:MAX(screenHeightPixels, screenWidthPixels)];
    [ac registerStaticFloatProperty:@"screen_scale"
                              value:UIScreen.mainScreen.scale];

    [CMDefaultProperties setUserInterfaceIdiom];

    [CMDefaultProperties setDeviceModel];

    // Battery
    CMBatteryLevelPropertyProvider *batteryLevelProvider =
        [[CMBatteryLevelPropertyProvider alloc] init];
    [ac registerLibPropertyProvider:@"device_battery_level"
                                dpp:batteryLevelProvider];
    CMBatteryStatePropertyProvider *batteryStateProvider =
        [[CMBatteryStatePropertyProvider alloc] init];
    [ac registerLibPropertyProvider:@"device_battery_state"
                                dpp:batteryStateProvider];
    CMLowPowerModePropertyProvider *lowPowerModeProvider =
        [[CMLowPowerModePropertyProvider alloc] init];
    [ac registerLibPropertyProvider:@"device_low_power_mode"
                                dpp:lowPowerModeProvider];
}

+ (void)setVersionString:(NSString *)versionString
               forPrefix:(NSString *)prefix {
    NSError *error;
    [AppcoreSharedAppcore() registerStaticVersionNumberProperty:prefix
                                                  versionString:versionString
                                                          error:&error];
    if (error) {
        NSLog(
            @"CriticalMoments: issue saving version number property: %@ -- %@",
            prefix, versionString);
    }
}

+ (void)setUserInterfaceIdiom {

    NSString *stringUserInterfaceIdiom = @"unknown";
    switch (UIDevice.currentDevice.userInterfaceIdiom) {
    case UIUserInterfaceIdiomPhone:
        stringUserInterfaceIdiom = @"phone";
        break;
    case UIUserInterfaceIdiomPad:
        stringUserInterfaceIdiom = @"tablet";
        break;
    case UIUserInterfaceIdiomTV:
        stringUserInterfaceIdiom = @"tv";
        break;
    case UIUserInterfaceIdiomCarPlay:
        stringUserInterfaceIdiom = @"car";
        break;
    case UIUserInterfaceIdiomMac:
        stringUserInterfaceIdiom = @"computer";
        break;

    default:
        break;
    }
    [AppcoreSharedAppcore()
        registerStaticStringProperty:@"user_interface_idiom"
                               value:stringUserInterfaceIdiom];
}

+ (void)setDeviceModel {
    AppcoreAppcore *ac = AppcoreSharedAppcore();

    [ac registerStaticStringProperty:@"device_manufacturer" value:@"Apple"];
    [ac registerStaticStringProperty:@"device_model_class"
                               value:UIDevice.currentDevice.model];

    struct utsname systemInfo;
    uname(&systemInfo);

    NSString *deviceModel = [NSString stringWithCString:systemInfo.machine
                                               encoding:NSUTF8StringEncoding];

    if (deviceModel == nil || deviceModel.length == 0) {
        [ac registerStaticStringProperty:@"device_model" value:@"unknown"];
        return;
    }

    if ([@[ @"arm64", @"i386", @"x86_64" ] containsObject:deviceModel]) {
        // This is a simulator. They don't return a model_version_number
        [ac registerStaticStringProperty:@"device_model" value:@"simulator"];
        return;
    }

    // format:
    // https://everyi.com/by-identifier/ipod-iphone-ipad-specs-by-model-identifier.html
    [ac registerStaticStringProperty:@"device_model" value:deviceModel];
    // remove non numeric chars, and replace comma with .
    NSString *numericString = [[deviceModel
        componentsSeparatedByCharactersInSet:
            [[NSCharacterSet characterSetWithCharactersInString:@"0123456789,."]
                invertedSet]] componentsJoinedByString:@""];
    numericString = [numericString stringByReplacingOccurrencesOfString:@","
                                                             withString:@"."];
    if (numericString.length > 0) {
        [CMDefaultProperties setVersionString:numericString
                                    forPrefix:@"device_model"];
    }
}

@end
