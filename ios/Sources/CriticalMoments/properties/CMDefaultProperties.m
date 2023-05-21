//
//  CMDefaultProperties.m
//
//
//  Created by Steve Cosman on 2023-05-20.
//

#import "CMDefaultProperties.h"

@import UIKit;

@implementation CMDefaultProperties

+ (void)registerDefaultPropertiesToAppcore {
    AppcoreAppcore *ac = AppcoreSharedAppcore();
    [ac registerStaticStringProperty:@"platform"
                               value:UIDevice.currentDevice.systemName];
    [ac registerStaticStringProperty:@"os_version_string"
                               value:UIDevice.currentDevice.systemVersion];
    [ac registerStaticStringProperty:@"device_manufacturer" value:@"Apple"];
    [ac registerStaticStringProperty:@"device_model"
                               value:UIDevice.currentDevice.model];

    NSLocale *locale = [NSLocale currentLocale];
    [ac registerStaticStringProperty:@"locale_language_code"
                               value:locale.languageCode];
    [ac registerStaticStringProperty:@"locale_country_code"
                               value:locale.countryCode];
    [ac registerStaticStringProperty:@"locale_currency_code"
                               value:locale.currencyCode];

    [ac registerStaticStringProperty:@"app_id"
                               value:NSBundle.mainBundle.bundleIdentifier];
    NSString *appVersion = [NSBundle.mainBundle
        objectForInfoDictionaryKey:@"CFBundleShortVersionString"];
    [ac registerStaticStringProperty:@"app_version_string" value:appVersion];

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
    [ac registerStaticStringProperty:@"user_interface_idiom"
                               value:stringUserInterfaceIdiom];
}

@end
