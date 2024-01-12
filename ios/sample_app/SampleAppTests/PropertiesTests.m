//
//  PropertiesTests.m
//  SampleAppTests
//
//  Created by Steve Cosman on 2023-10-25.
//

#import <XCTest/XCTest.h>

#import "../SampleApp/AppDelegate.h"

@interface PropertiesTests : XCTestCase

@end

@implementation PropertiesTests

- (void)setUp {
}

- (void)tearDown {
}

- (void)testPropertyProviders {
    // This test iterates all property providers, and tests they provide valid values
    // E2E test because it uses entire stack (App integration, SPM, Appcore, SDK, SDK property providers, conditions)
    // Running on device is better than simulator as some properties have more rich data on device (battery level,
    // networking, orientation, etc), but should pass on any platform/device.

    // Important: these tests should reflect the exact values from our documentation as a confirmation docs are correct
    // For example: only allow nil if we have documented that the value is nullable

    // clang-format off
    NSDictionary *cases = @{
        @"platform" : @"platform in ['iOS', 'iPadOS']",
        @"user_interface_idiom": @"user_interface_idiom in ['phone', 'tablet', 'tv', 'car', 'computer', 'unknown']", // add_test_count
        @"os_version": @"os_version != nil && versionGreaterThan(os_version, '15.0') && versionGreaterThan('99.0', os_version)", // add_test_count
        @"cm_version": @"cm_version != nil && versionGreaterThan('999.0', cm_version)", // add_test_count
        @"device_manufacturer": @"device_manufacturer == 'Apple'", // add_test_count
        @"device_model_class": @"device_model_class in ['iPhone', 'iPod', 'iPad']", // add_test_count
        @"device_model": @"device_model != nil", // add_test_count
        // device_model_version only set on devices, not simulators
        @"device_model_version": @"device_model == 'simulator' || (device_model_version != nil && versionGreaterThan('999.0', device_model_version))" , // add_test_count
        @"dark_mode": @"dark_mode in [true, false]", // add_test_count
        @"interface_orientation": @"interface_orientation in ['landscape', 'portrait']", // add_test_count
        @"device_orientation": @"device_orientation in ['landscape', 'portrait', 'face_up', 'face_down', 'unknown']", // add_test_count
        @"screen_width_points": @"screen_width_points != nil && screen_width_points > 0 && screen_width_points < 99999", // add_test_count
        @"screen_height_points": @"screen_height_points != nil && screen_height_points > 0 && screen_height_points < 99999", // add_test_count
        @"screen_width_pixels": @"screen_width_pixels != nil && screen_width_pixels > 0 && screen_width_pixels < 99999", // add_test_count
        @"screen_height_pixels": @"screen_height_pixels != nil && screen_height_pixels > 0 && screen_height_pixels < 99999", // add_test_count
        @"screen_scale": @"screen_scale != nil && screen_scale >= 1.0 && screen_scale <= 10.0", // add_test_count
        @"screenScalePixToPoint": @"screen_width_pixels / screen_scale == screen_width_points && screen_height_pixels / screen_scale == screen_height_points", // add_test_count
        @"screen_brightness": @"screen_brightness >= 0.0 && screen_brightness <= 1.0", // add_test_count
        @"screen_captured": @"screen_captured in [true,false]", // add_test_count
        @"locale_language_code": @"locale_language_code != nil && len(locale_language_code) == 2", // add_test_count
        @"locale_country_code": @"locale_country_code != nil && len(locale_country_code) == 2", // add_test_count
        @"locale_currency_code": @"locale_currency_code != nil && len(locale_currency_code) == 3", // add_test_count
        @"locale_language_direction": @"locale_language_direction in ['RTL', 'LTR']", // add_test_count
        @"app_id": @"app_id == 'io.criticalmoments.sample-app'", // add_test_count
        @"app_version": @"app_version == '1.0'", // add_test_count
        @"app_install_date": @"app_install_date != nil && app_install_date > unixTimeMilliseconds(1688744356123) && app_install_date < unixTimeMilliseconds(1988744356123)", // add_test_count
        @"app_install_date_now": @"app_install_date <= now()", // add_test_count
        @"device_battery_level": @"device_battery_level == -1 || (device_battery_level >= 0.0 && device_battery_level <= 1.0)", // add_test_count
        @"device_battery_state": @"device_battery_state in ['charging', 'full', 'unplugged', 'unknown']", // add_test_count
        @"device_low_power_mode": @"device_low_power_mode in [true,false]", // add_test_count
        @"network_connection_type": @"network_connection_type in ['wifi', 'cellular', 'wired', 'unknown']", // add_test_count
        @"has_active_network": @"has_active_network in [true,false]", // add_test_count
        @"low_data_mode": @"low_data_mode in [true,false,nil]", // add_test_count
        @"expensive_network": @"expensive_network in [true, false]", // add_test_count
        @"has_wifi_connection": @"has_wifi_connection in [true, false]", // add_test_count
        @"has_cell_connection": @"has_cell_connection in [true, false]", // add_test_count
        @"on_call": @"on_call in [true, false]", // add_test_count
        @"other_audio_playing": @"other_audio_playing in [true, false]", // add_test_count
        @"has_watch": @"has_watch in [true, false]", // add_test_count
        @"foreground": @"foreground in [true, false]", // add_test_count
        @"app_state": @"app_state in ['active', 'inactive', 'background', 'unknown']", // add_test_count
        @"session_start_time": @"now() > session_start_time && now() - duration('5m') < session_start_time", // add_test_count
        
        // Audio
        @"has_headphones": @"has_headphones in [true,false]", // add_test_count
        @"has_bt_headphones": @"has_bt_headphones in [true,false]", // add_test_count
        @"has_bt_headset": @"has_bt_headset in [true,false]", // add_test_count
        @"has_wired_headset": @"has_wired_headset in [true,false]", // add_test_count
        @"has_car_audio": @"has_car_audio in [true,false]", // add_test_count

        @"rand": @"(rand() % 100) >= 0 && (rand() % 100) < 100", // add_test_count
        @"sessionRand": @"(sessionRand() % 100) >= 0 && (sessionRand() % 100) < 100 && sessionRand() == sessionRand()", // add_test_count
        @"randForKey": @"randForKey('key1', 1) == 292785326893130985", // add_test_count
        
        @"timezone_gmt_offset": @"timezone_gmt_offset != nil && timezone_gmt_offset <= 24*60*60 && timezone_gmt_offset >= -24*60*60", // add_test_count
        @"location_permission": @"location_permission in [true,false]", // add_test_count
        @"location_permission_detailed": @"location_permission_detailed in ['not_determined', 'restricted', 'denied', 'authorized_always', 'authorized_when_in_use', 'unknown']", // add_test_count
        // TODO P2: simulate location and check actual values
        @"location_latitude": @"(location_latitude ?? 0) <= 90.0 && (location_latitude ?? 0) >= -90.0", // add_test_count
        @"location_longitude": @"(location_longitude ?? 0) <= 180.0 && (location_longitude ?? 0) >= -180.0", // add_test_count
        @"location_city": @"location_city == nil || len (location_city ?? '') > 0", // add_test_count
        @"location_region": @"location_region == nil || len (location_region ?? '') > 0", // add_test_count
        @"location_country": @"location_country == nil || len (location_country ?? '') > 0", // add_test_count
        
        // approx location will change based on ip
        @"location_approx_latitude": @"(location_approx_latitude ?? 0) <= 90.0 && (location_approx_latitude ?? 0) >= -90.0", // add_test_count
        @"location_approx_longitude": @"(location_approx_longitude ?? 0) <= 180.0 && (location_approx_longitude ?? 0) >= -180.0", // add_test_count
        @"location_approx_city": @"location_approx_city == nil || len (location_approx_city ?? '') > 0", // add_test_count
        @"location_approx_region": @"location_approx_region == nil || len (location_approx_region ?? '') > 0", // add_test_count
        @"location_approx_country": @"location_approx_country == nil || len (location_approx_country ?? '') > 0", // add_test_count
        
        // Weather -- tested in library
        
        @"contacts_permission": @"contacts_permission in ['not_determined', 'restricted', 'denied', 'authorized', 'unknown']", // add_test_count
        @"camera_permission": @"camera_permission in ['not_determined', 'restricted', 'denied', 'authorized', 'unknown']", // add_test_count
        @"microphone_permission": @"microphone_permission in ['not_determined', 'restricted', 'denied', 'authorized', 'unknown']", // add_test_count
        @"notifications_permission": @"notifications_permission in ['not_determined', 'denied', 'authorized', 'provisional', 'ephemeral', 'unknown']", // add_test_count
        @"photo_library_permission": @"photo_library_permission in ['not_determined', 'denied', 'authorized', 'restricted', 'limited', 'unknown']", // add_test_count
        @"add_photo_permission": @"add_photo_permission in ['not_determined', 'denied', 'authorized', 'restricted', 'limited', 'unknown']", // add_test_count
        @"calendar_permission": @"calendar_permission in ['not_determined', 'denied', 'authorized_full', 'authorized_write_only', 'restricted', 'unknown']", // add_test_count
        @"reminders_permission": @"reminders_permission in ['not_determined', 'denied', 'authorized_full', 'authorized_write_only', 'restricted', 'unknown']", // add_test_count
        @"bluetooth_permission": @"bluetooth_permission in ['not_determined', 'restricted', 'denied', 'authorized', 'unknown']", // add_test_count
        
        // Functions
        @"propertyHistoryLatestValue": @"propertyHistoryLatestValue('platform') == 'iOS' || propertyHistoryLatestValue('platform') == 'iPadOS'", // add_test_count
        @"propertyHistoryLatestValueNil": @"propertyHistoryLatestValue('never_set_prop') == nil", // add_test_count
        @"propertyEver": @"propertyEver('app_id', 'io.criticalmoments.sample-app') && !propertyEver('app_id', 'wrongval') && !propertyEver('wrongproperty', 'a')", // add_test_count
        @"stableRand": @"stableRand() == stableRand()", // add_test_count
    };
    // clang-format on

    id<UIApplicationDelegate> ad = UIApplication.sharedApplication.delegate;
    AppDelegate *aad = (AppDelegate *)ad;
    CriticalMoments *cm = [aad cmInstance];

    // Expectations just used to test condition evaluated, not that it passed
    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    for (NSString *name in cases.keyEnumerator) {
        NSString *condition = cases[name];

        // Expectations are only used to wait -- actual assets in the callback
        XCTestExpectation *expectation = [[XCTestExpectation alloc] initWithDescription:name];
        [expectations addObject:expectation];
        [cm checkNamedCondition:name
                      condition:condition
                        handler:^(bool result, NSError *_Nullable error) {
                          if (error != nil) {
                              XCTAssert(false, @"Property test failed with error: %@", error);
                          }
                          XCTAssertTrue(result, @"Property test did not return true for condition check: %@", name);
                          [expectation fulfill];
                        }];
    }

    [self waitForExpectations:expectations timeout:20];
}

- (void)testGeoIpLocation {
    id<UIApplicationDelegate> ad = UIApplication.sharedApplication.delegate;
    AppDelegate *aad = (AppDelegate *)ad;
    CriticalMoments *cm = [aad cmInstance];

    NSString *condition =
        @"(location_approx_city == nil || len(location_approx_city) > 0) && (location_approx_country == nil || "
        @"len(location_approx_country) > 0) && (location_approx_region == nil || "
        @"len(location_approx_region) > 0) && (location_approx_latitude == nil || abs(location_approx_latitude) <= 90) "
        @"&& (location_approx_longitude == nil || abs(location_approx_longitude) <= 180)";

    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    XCTestExpectation *expectation1 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation1];
    [cm checkNamedCondition:@"locCondition"
                  condition:condition
                    handler:^(bool result, NSError *error) {
                      if (!result || error) {
                          XCTAssert(false, "approx location condition failed to return");
                      }
                      [expectation1 fulfill];
                    }];

    [self waitForExpectations:expectations timeout:5.0];
}

// Not included in test plan, so not run by default. But helpful for development.
- (void)testGeoIpLocationToronto {
    id<UIApplicationDelegate> ad = UIApplication.sharedApplication.delegate;
    AppDelegate *aad = (AppDelegate *)ad;
    CriticalMoments *cm = [aad cmInstance];

    NSString *condition =
        @"location_approx_city == 'Toronto' && location_approx_country == 'CA' && location_approx_region == 'ON' "
        @"&& abs(location_approx_latitude - 43.651070) < 0.5 && abs(location_approx_longitude - -79.347015) < 0.5";

    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    XCTestExpectation *expectation1 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation1];
    [cm checkNamedCondition:@"locCondition"
                  condition:condition
                    handler:^(bool result, NSError *error) {
                      if (!result || error) {
                          XCTAssert(false, "approx location condition failed to return");
                      }
                      [expectation1 fulfill];
                    }];

    [self waitForExpectations:expectations timeout:5.0];
}

@end
