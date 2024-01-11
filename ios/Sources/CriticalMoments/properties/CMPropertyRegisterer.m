//
//  CMDefaultProperties.m
//
//
//  Created by Steve Cosman on 2023-05-20.
//

#import "CMPropertyRegisterer.h"

#import "../CriticalMoments_private.h"
#import "../utils/CMUtils.h"
#import "CMAudioPropertyProvider.h"
#import "CMBatteryLevelPropertyProvider.h"
#import "CMCallPropertyProvider.h"
#import "CMLocationPropertyProvider.h"
#import "CMMiscPropertyProviders.h"
#import "CMNetworkingPropertyProvider.h"
#import "CMPermissionsPropertyProvider.h"
#import "CMViewPropertyProvider.h"

#import <sys/utsname.h>

@import UIKit;

@interface CMPropertyRegisterer ()

@property(nonatomic, strong) AppcoreAppcore *appcore;

@end

@implementation CMPropertyRegisterer

- (instancetype)initWithAppcore:(AppcoreAppcore *)ac {
    self = [super init];
    if (self) {
        _appcore = ac;
    }
    return self;
}

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
    [_appcore registerStaticStringProperty:key value:value error:&error];
    [self processError:error];
}

- (void)registerStaticIntProperty:(NSString *)key value:(long)value {
    NSError *error;
    [_appcore registerStaticIntProperty:key value:value error:&error];
    [self processError:error];
}

- (void)registerStaticFloatProperty:(NSString *)key value:(double)value {
    NSError *error;
    [_appcore registerStaticFloatProperty:key value:value error:&error];
    [self processError:error];
}

- (void)registerStaticBoolProperty:(NSString *)key value:(bool)value {
    NSError *error;
    [_appcore registerStaticBoolProperty:key value:value error:&error];
    [self processError:error];
}

- (void)registerStaticTimeProperty:(NSString *)key value:(NSDate *)value {
    NSError *error;
    int64_t goTime = [CMUtils dateToGoTime:value];
    [_appcore registerStaticTimeProperty:key value:goTime error:&error];
    [self processError:error];
}

- (void)registerLibPropertyProvider:(NSString *)key value:(id<CMDynamicPropertyProvider>)value {
    NSError *error;
    // Wrap the CMDynamicPropertyProvider to implement the appcore interface
    CMDynamicPropertyProviderWrapper *wrapper = [[CMDynamicPropertyProviderWrapper alloc] initWithPP:value];
    [_appcore registerLibPropertyProvider:key dpp:wrapper error:&error];
    [self processError:error];
}

- (void)registerDefaultPropertiesToAppcore {

    // This API returns different values on older iPads. Make these
    // consistent (documented in docs)
    NSString *systemOsName = UIDevice.currentDevice.systemName;
    if ([@"iOS" isEqualToString:systemOsName] && UIDevice.currentDevice.userInterfaceIdiom == UIUserInterfaceIdiomPad) {
        systemOsName = @"iPadOS";
    }
    [self registerStaticStringProperty:@"platform" value:systemOsName];

    // OS Version
    [self registerStaticStringProperty:@"os_version" value:UIDevice.currentDevice.systemVersion];

    // Locale
    NSLocale *locale = [NSLocale currentLocale];
    [self registerStaticStringProperty:@"locale_language_code" value:locale.languageCode];
    [self registerStaticStringProperty:@"locale_country_code" value:locale.countryCode];
    [self registerStaticStringProperty:@"locale_currency_code" value:locale.currencyCode];
    CMLanguageDirectionPropertyProvider *ldpp = [[CMLanguageDirectionPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"locale_language_direction" value:ldpp];

    // Session start time
    NSDate *now = [[NSDate alloc] init];
    [self registerStaticTimeProperty:@"session_start_time" value:now];

    // Bundle ID
    [self registerStaticStringProperty:@"app_id" value:NSBundle.mainBundle.bundleIdentifier];

    // App Version
    NSString *appVersion = [NSBundle.mainBundle objectForInfoDictionaryKey:@"CFBundleShortVersionString"];
    [self registerStaticStringProperty:@"app_version" value:appVersion];

    // Library Version
    [self registerStaticStringProperty:@"cm_version" value:CM_LIB_VERSION_NUMBER_STRING];

    // Screen size / scale
    CGSize screenSize = UIScreen.mainScreen.bounds.size;
    [self registerStaticIntProperty:@"screen_width_points" value:MIN(screenSize.width, screenSize.height)];
    [self registerStaticIntProperty:@"screen_height_points" value:MAX(screenSize.width, screenSize.height)];
    CGFloat screenWidthPixels = screenSize.width * UIScreen.mainScreen.scale;
    CGFloat screenHeightPixels = screenSize.height * UIScreen.mainScreen.scale;
    [self registerStaticIntProperty:@"screen_width_pixels" value:MIN(screenHeightPixels, screenWidthPixels)];
    [self registerStaticIntProperty:@"screen_height_pixels" value:MAX(screenHeightPixels, screenWidthPixels)];
    [self registerStaticFloatProperty:@"screen_scale" value:UIScreen.mainScreen.scale];
    CMBrightnessProvider *brpp = [[CMBrightnessProvider alloc] init];
    [self registerLibPropertyProvider:@"screen_brightness" value:brpp];
    CMScreenCapturedProvider *scpp = [[CMScreenCapturedProvider alloc] init];
    [self registerLibPropertyProvider:@"screen_captured" value:scpp];

    [self setUserInterfaceIdiom];

    [self setDeviceModel];

    // Battery
    CMBatteryLevelPropertyProvider *batteryLevelProvider = [[CMBatteryLevelPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"device_battery_level" value:batteryLevelProvider];
    CMBatteryStatePropertyProvider *batteryStateProvider = [[CMBatteryStatePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"device_battery_state" value:batteryStateProvider];
    CMLowPowerModePropertyProvider *lowPowerModeProvider = [[CMLowPowerModePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"device_low_power_mode" value:lowPowerModeProvider];

    // Screen/views
    CMDeviceOrientationPropertyProvider *deviceOrientationProvider = [[CMDeviceOrientationPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"device_orientation" value:deviceOrientationProvider];
    CMInterfaceOrientationPropertyProvider *uiOrientationProvider =
        [[CMInterfaceOrientationPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"interface_orientation" value:uiOrientationProvider];
    CMDarkModePropertyProvider *darkModeProvider = [[CMDarkModePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"dark_mode" value:darkModeProvider];
    CMAppStatePropertyProvider *appStateProvider = [[CMAppStatePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"app_state" value:appStateProvider];
    CMForegroundProvider *foregroundProvider = [[CMForegroundProvider alloc] init];
    [self registerLibPropertyProvider:@"foreground" value:foregroundProvider];

    // Networking

    CMHasActiveNetworkPropertyProvider *hasActiveNetworkPP = [[CMHasActiveNetworkPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"has_active_network" value:hasActiveNetworkPP];

    CMNetworkTypePropertyProvider *networkTypePP = [[CMNetworkTypePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"network_connection_type" value:networkTypePP];

    // low network mode added in ios 13
    BOOL deviceHasLowDataMode = false;
    if (@available(iOS 13, *)) {
        deviceHasLowDataMode = true;
    }
    if (deviceHasLowDataMode) {
        CMLowDataModePropertyProvider *lowDataProvider = [[CMLowDataModePropertyProvider alloc] init];
        [self registerLibPropertyProvider:@"low_data_mode" value:lowDataProvider];
    }

    CMExpensiveNetworkPropertyProvider *expensiveNetworkPP = [[CMExpensiveNetworkPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"expensive_network" value:expensiveNetworkPP];

    CMHasWifiConnectionPropertyProvider *hasWifiPP = [[CMHasWifiConnectionPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"has_wifi_connection" value:hasWifiPP];

    CMHasCellConnectionPropertyProvider *hasCellPP = [[CMHasCellConnectionPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"has_cell_connection" value:hasCellPP];

    // Audio
    CMAudioPlayingPropertyProvider *audioPlayingProvider = [[CMAudioPlayingPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"other_audio_playing" value:audioPlayingProvider];
    [self registerLibPropertyProvider:@"has_headphones" value:[CMAudioPortPropertyProvider hasHeadphones]];
    [self registerLibPropertyProvider:@"has_bt_headphones" value:[CMAudioPortPropertyProvider hasBtHeadphones]];
    [self registerLibPropertyProvider:@"has_bt_headset" value:[CMAudioPortPropertyProvider hasBtHeadset]];
    [self registerLibPropertyProvider:@"has_wired_headset" value:[CMAudioPortPropertyProvider hasWiredHeadset]];
    [self registerLibPropertyProvider:@"has_car_audio" value:[CMAudioPortPropertyProvider hasCarAudio]];

    // Calls
    CMCallPropertyProvider *callsPP = [[CMCallPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"on_call" value:callsPP];

    // Misc
    CMAppInstallDatePropertyProviders *appInstallProvider = [[CMAppInstallDatePropertyProviders alloc] init];
    [self registerLibPropertyProvider:@"app_install_date" value:appInstallProvider];
    CMHasWatchPropertyProviders *hasWatchProvider = [[CMHasWatchPropertyProviders alloc] init];
    [self registerLibPropertyProvider:@"has_watch" value:hasWatchProvider];

    // Location
    CMLocationPermissionsPropertyProvider *lppp = [[CMLocationPermissionsPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_permission" value:lppp];
    CMLocationPermissionDetailedPropertyProvider *lpdpp = [[CMLocationPermissionDetailedPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_permission_detailed" value:lpdpp];
    CMLatitudePropertyProvider *latpp = [[CMLatitudePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_latitude" value:latpp];
    CMLongitudePropertyProvider *longpp = [[CMLongitudePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_longitude" value:longpp];
    CMCityPropertyProvider *citypp = [[CMCityPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_city" value:citypp];
    CMRegionPropertyProvider *regionpp = [[CMRegionPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_region" value:regionpp];
    CMCountryPropertyProvider *countrypp = [[CMCountryPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_country" value:countrypp];

    // Approx Location
    CMApproxCityPropertyProvider *approxCity = [[CMApproxCityPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_approx_city" value:approxCity];
    CMApproxRegionPropertyProvider *approxRegion = [[CMApproxRegionPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_approx_region" value:approxRegion];
    CMApproxCountryPropertyProvider *approxCountry = [[CMApproxCountryPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_approx_country" value:approxCountry];
    CMApproxLatitudePropertyProvider *approxLat = [[CMApproxLatitudePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_approx_latitude" value:approxLat];
    CMApproxLongitudePropertyProvider *approxLong = [[CMApproxLongitudePropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"location_approx_longitude" value:approxLong];

    // Permissions
    CMNotificationPermissionsPropertyProvider *npp = [[CMNotificationPermissionsPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"notifications_permission" value:npp];
    CMCapturePermissionsPropertyProvider *micpp =
        [[CMCapturePermissionsPropertyProvider alloc] initWithMediaType:AVMediaTypeAudio];
    [self registerLibPropertyProvider:@"microphone_permission" value:micpp];
    CMCapturePermissionsPropertyProvider *campp =
        [[CMCapturePermissionsPropertyProvider alloc] initWithMediaType:AVMediaTypeVideo];
    [self registerLibPropertyProvider:@"camera_permission" value:campp];
    CMContactsPermissionsPropertyProvider *conpp = [[CMContactsPermissionsPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"contacts_permission" value:conpp];
    if (@available(iOS 14.0, *)) {
        CMPhotosPermissionsPropertyProvider *plpp =
            [[CMPhotosPermissionsPropertyProvider alloc] initWithAccessLevel:PHAccessLevelReadWrite];
        [self registerLibPropertyProvider:@"photo_library_permission" value:plpp];
        CMPhotosPermissionsPropertyProvider *appp =
            [[CMPhotosPermissionsPropertyProvider alloc] initWithAccessLevel:PHAccessLevelAddOnly];
        [self registerLibPropertyProvider:@"add_photo_permission" value:appp];
    } else {
        // same permission for both prior to iOS 14
        CMPhotosPermissionsPropertyProvider *ppp = [[CMPhotosPermissionsPropertyProvider alloc] init];
        [self registerLibPropertyProvider:@"photo_library_permission" value:ppp];
        [self registerLibPropertyProvider:@"add_photo_permission" value:ppp];
    }
    CMCalendarPermissionsPropertyProvider *calpp =
        [[CMCalendarPermissionsPropertyProvider alloc] initWithEntityType:EKEntityTypeEvent];
    [self registerLibPropertyProvider:@"calendar_permission" value:calpp];
    CMCalendarPermissionsPropertyProvider *rempp =
        [[CMCalendarPermissionsPropertyProvider alloc] initWithEntityType:EKEntityTypeReminder];
    [self registerLibPropertyProvider:@"reminders_permission" value:rempp];
    CMBluetoothPermissionsPropertyProvider *btpp = [[CMBluetoothPermissionsPropertyProvider alloc] init];
    [self registerLibPropertyProvider:@"bluetooth_permission" value:btpp];

    // Weather
    if (@available(iOS 16.0, *)) {
        NSDictionary<NSString *, CMWeatherPropertyProvider *> *weatherProviders =
            [CMWeatherPropertyProvider allWeatherProviders];
        for (NSString *conditionName in weatherProviders.keyEnumerator) {
            CMWeatherPropertyProvider *provider = weatherProviders[conditionName];
            [self registerLibPropertyProvider:conditionName value:provider];
        }
    }
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
    [self registerStaticStringProperty:@"user_interface_idiom" value:stringUserInterfaceIdiom];
}

- (void)setDeviceModel {
    [self registerStaticStringProperty:@"device_manufacturer" value:@"Apple"];
    [self registerStaticStringProperty:@"device_model_class" value:UIDevice.currentDevice.model];

    struct utsname systemInfo;
    uname(&systemInfo);

    NSString *deviceModel = [NSString stringWithCString:systemInfo.machine encoding:NSUTF8StringEncoding];

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
        componentsSeparatedByCharactersInSet:[[NSCharacterSet characterSetWithCharactersInString:@"0123456789,."]
                                                 invertedSet]] componentsJoinedByString:@""];
    versionString = [versionString stringByReplacingOccurrencesOfString:@"," withString:@"."];
    if (versionString.length > 0) {
        [self registerStaticStringProperty:@"device_model_version" value:versionString];
    }
}

@end
