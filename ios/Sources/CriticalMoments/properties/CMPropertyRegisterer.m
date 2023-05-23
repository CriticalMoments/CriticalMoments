//
//  CMDefaultProperties.m
//
//
//  Created by Steve Cosman on 2023-05-20.
//

#import "CMPropertyRegisterer.h"

#import "CMBatteryLevelPropertyProvider.h"

#import <sys/utsname.h>

@import UIKit;

@implementation CMPropertyRegisterer

- (void)processError:(NSError *)error {
    if (!error) {
        return;
    }
    NSLog(@"CriticalMoments: Issue registering properties \"%@\"", error);
#if DEBUG
    NSLog(@"CriticalMoments: throwing a "
          @"NSInternalInconsistencyException "
          @"to help find this issue. Exceptions are only thrown in "
          @"debug "
          @"mode, and will not crash apps built for release.");
    @throw NSInternalInconsistencyException;
#endif
}

- (void)registerStaticStringProperty:(NSString *)key value:(NSString *)value {
    NSError *error;
    [AppcoreSharedAppcore() registerStaticStringProperty:key
                                                   value:value
                                                   error:&error];
    [self processError:error];
}

- (void)registerStaticIntProperty:(NSString *)key value:(long)value {
    NSError *error;
    [AppcoreSharedAppcore() registerStaticIntProperty:key
                                                value:value
                                                error:&error];
    [self processError:error];
}

- (void)registerStaticFloatProperty:(NSString *)key value:(double)value {
    NSError *error;
    [AppcoreSharedAppcore() registerStaticFloatProperty:key
                                                  value:value
                                                  error:&error];
    [self processError:error];
}

- (void)registerStaticBoolProperty:(NSString *)key value:(bool)value {
    NSError *error;
    [AppcoreSharedAppcore() registerStaticBoolProperty:key
                                                 value:value
                                                 error:&error];
    [self processError:error];
}

- (void)registerLibPropertyProvider:(NSString *)key
                              value:(id<AppcoreLibPropertyProvider>)value {
    NSError *error;
    [AppcoreSharedAppcore() registerLibPropertyProvider:key
                                                    dpp:value
                                                  error:&error];
    [self processError:error];
}

- (void)registerStaticVersionProperty:(NSString *)key value:(NSString *)value {
    NSError *error;
    [AppcoreSharedAppcore() registerStaticVersionNumberProperty:key
                                                  versionString:value
                                                          error:&error];
    [self processError:error];
}

- (void)registerDefaultPropertiesToAppcore {

    // This API returns different values on older iPads. Make these
    // consistent (documented in docs)
    NSString *systemOsName = UIDevice.currentDevice.systemName;
    if ([@"iOS" isEqualToString:systemOsName] &&
        UIDevice.currentDevice.userInterfaceIdiom == UIUserInterfaceIdiomPad) {
        systemOsName = @"iPadOS";
    }
    [self registerStaticStringProperty:@"platform" value:systemOsName];

    // OS Version
    [self registerStaticVersionProperty:@"os_version"
                                  value:UIDevice.currentDevice.systemVersion];

    // Locale
    NSLocale *locale = [NSLocale currentLocale];
    [self registerStaticStringProperty:@"locale_language_code"
                                 value:locale.languageCode];
    [self registerStaticStringProperty:@"locale_country_code"
                                 value:locale.countryCode];
    [self registerStaticStringProperty:@"locale_currency_code"
                                 value:locale.currencyCode];

    // Bundle ID
    [self registerStaticStringProperty:@"app_id"
                                 value:NSBundle.mainBundle.bundleIdentifier];

    // App Version
    NSString *appVersion = [NSBundle.mainBundle
        objectForInfoDictionaryKey:@"CFBundleShortVersionString"];
    [self registerStaticVersionProperty:@"app_version" value:appVersion];

    // Screen size / scale
    CGSize screenSize = UIScreen.mainScreen.bounds.size;
    [self registerStaticIntProperty:@"screen_width_points"
                              value:MIN(screenSize.width, screenSize.height)];
    [self registerStaticIntProperty:@"screen_height_points"
                              value:MAX(screenSize.width, screenSize.height)];
    CGFloat screenWidthPixels = screenSize.width * UIScreen.mainScreen.scale;
    CGFloat screenHeightPixels = screenSize.height * UIScreen.mainScreen.scale;
    [self registerStaticIntProperty:@"screen_width_pixels"
                              value:MIN(screenHeightPixels, screenWidthPixels)];
    [self registerStaticIntProperty:@"screen_height_pixels"
                              value:MAX(screenHeightPixels, screenWidthPixels)];
    [self registerStaticFloatProperty:@"screen_scale"
                                value:UIScreen.mainScreen.scale];

    [self setUserInterfaceIdiom];

    [self setDeviceModel];

    // Battery
    CMBatteryLevelPropertyProvider *batteryLevelProvider =
        [[CMBatteryLevelPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"device_battery_level"
                                value:batteryLevelProvider];
    CMBatteryStatePropertyProvider *batteryStateProvider =
        [[CMBatteryStatePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"device_battery_state"
                                value:batteryStateProvider];
    CMLowPowerModePropertyProvider *lowPowerModeProvider =
        [[CMLowPowerModePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"device_low_power_mode"
                                value:lowPowerModeProvider];
}

- (void)setUserInterfaceIdiom {

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
    [self registerStaticStringProperty:@"user_interface_idiom"
                                 value:stringUserInterfaceIdiom];
}

- (void)setDeviceModel {
    [self registerStaticStringProperty:@"device_manufacturer" value:@"Apple"];
    [self registerStaticStringProperty:@"device_model_class"
                                 value:UIDevice.currentDevice.model];

    struct utsname systemInfo;
    uname(&systemInfo);

    NSString *deviceModel = [NSString stringWithCString:systemInfo.machine
                                               encoding:NSUTF8StringEncoding];

    if (deviceModel == nil || deviceModel.length == 0) {
        [self registerStaticStringProperty:@"device_model" value:@"unknown"];
        return;
    }

    if ([@[ @"arm64", @"i386", @"x86_64" ] containsObject:deviceModel]) {
        // This is a simulator. They don't return a model_version_number
        [self registerStaticStringProperty:@"device_model" value:@"simulator"];
        return;
    }

    // format:
    // https://everyi.com/by-identifier/ipod-iphone-ipad-specs-by-model-identifier.html
    [self registerStaticStringProperty:@"device_model" value:deviceModel];
    // remove non numeric chars, and replace comma with .
    NSString *versionString = [[deviceModel
        componentsSeparatedByCharactersInSet:
            [[NSCharacterSet characterSetWithCharactersInString:@"0123456789,."]
                invertedSet]] componentsJoinedByString:@""];
    versionString = [versionString stringByReplacingOccurrencesOfString:@","
                                                             withString:@"."];
    if (versionString.length > 0) {
        [self registerStaticVersionProperty:@"device_model_version"
                                      value:versionString];
    }
}

@end