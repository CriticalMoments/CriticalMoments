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

    // This API returns different values on older iOS. Make these consistent.
    // iOS, iPadOS (iPad and iPad app on Mac), tvOS
    NSString *systemOsName = UIDevice.currentDevice.systemName;
    if ([@"iOS" isEqualToString:systemOsName] &&
        UIDevice.currentDevice.userInterfaceIdiom == UIUserInterfaceIdiomPad) {
        systemOsName = @"iPadOS";
    }
    [ac registerStaticStringProperty:@"platform" value:systemOsName];

    NSError *error;
    [ac registerStaticVersionNumberProperty:@"os"
                              versionString:UIDevice.currentDevice.systemVersion
                                      error:&error];
    if (error) {
        NSLog(@"CriticalMoments: issue saving os version number property: %@",
              UIDevice.currentDevice.systemVersion);
    }

    // Make/Model
    [ac registerStaticStringProperty:@"device_manufacturer" value:@"Apple"];
    [ac registerStaticStringProperty:@"device_model"
                               value:UIDevice.currentDevice.model];

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
    error = nil;
    [ac registerStaticVersionNumberProperty:@"app"
                              versionString:appVersion
                                      error:&error];
    if (error) {
        NSLog(@"CriticalMoments: issue saving app version number property: %@",
              appVersion);
    }

    // UserInterfaceIdiom
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
