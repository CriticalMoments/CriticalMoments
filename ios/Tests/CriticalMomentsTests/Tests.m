//
//  CriticalMomentsTests.m
//  CriticalMomentsTests
//
//  Created by scosman on 04/12/2023.
//  Copyright (c) 2023 scosman. All rights reserved.
//

@import XCTest;

@import CriticalMoments;
@import Appcore;
#import "CriticalMoments.h"
#import "CriticalMoments_private.h"
@import EventKit;
@import Contacts;
#import "../../Sources/CriticalMoments/properties/CMLocationPropertyProvider.h"

// This key is only valid for test bundle "com.apple.dt.xctest.tool"
#define TEST_API_KEY                                                                                                   \
    @"CM1-Yjpjb20uYXBwbGUuZHQueGN0ZXN0LnRvb2w=-MEUCIQCktQU5A0wyr8rA7cHrrfZYxR/"                                        \
    @"7wTh+WIlgfLvOIeEFDQIgZGWxJeKshNah+0hP7J/5oH3V1CGvZyvAWrN+4WXfNoM="

@interface Tests : XCTestCase

@end

@implementation Tests

- (void)setUp {
    [super setUp];
    // Put setup code here. This method is called before the invocation of each
    // test method in the class.
}

- (void)tearDown {
    // Put teardown code here. This method is called after the invocation of
    // each test method in the class.
    [super tearDown];
}

- (void)testObjcPing {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];
    NSString *response = [cm objcPing];
    XCTAssert([@"objcPong" isEqualToString:response],
              @"Expected ping to pong -- objective C tests not working end to end");
}

- (void)testAppcoreIntegration {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];
    NSString *response = AppcoreGoPing();
    XCTAssert([@"AppcorePong->PongCmCore" isEqualToString:response],
              @"Expected ping to pong -- Appcore framework integration not "
              @"working end to end");

    NSString *fullyIntegratedRespons = [cm goPing];
    XCTAssert([@"AppcorePong->PongCmCore" isEqualToString:fullyIntegratedRespons],
              @"Expected ping to pong -- Appcore e2e framework integration not "
              @"working end to end");
}

- (void)testApiKeyValidation {

    NSError *error;
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];

    [cm setApiKey:@"" error:&error];
    XCTAssert(error != nil, @"Empty API key passed validation");

    error = nil;
    [cm setApiKey:@"invalid" error:&error];
    XCTAssert(error != nil, @"Invalid API key passed validation");

    error = nil;
    [cm setApiKey:@"CM1-aGVsbG86d29ybGQ=-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-"
                  @"MEUCIQCUfx6xlmQ0kdYkuw3SMFFI6WXrCWKWwetXBrXXG2hjAwIgWBPIMrdM1ET0Hbpn"
                  @"Xlnpj/f+VXtjRTqNNz9L/AOt4GY="
            error:&error];
    XCTAssert(error != nil, @"API key from another app passed validation");

    // Valid key
    error = nil;
    [cm setApiKey:TEST_API_KEY error:&error];
    XCTAssert(error == nil, @"API key failed validation");
}

- (CriticalMoments *)buildAndStartCMForTest {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];
    [self startCMForTest:cm];
    return cm;
}

- (void)startCMForTest:(CriticalMoments *)cm {
    [cm disableUserNotifications];
    NSBundle *testBundle = [NSBundle bundleForClass:self.class];
    NSURL *resourceBundleId =
        [testBundle.bundleURL URLByAppendingPathComponent:@"CriticalMoments_CriticalMomentsTests.bundle"];
    NSBundle *resourceBundle = [NSBundle bundleWithURL:resourceBundleId];
    NSURL *url = [resourceBundle URLForResource:@"TestResources/testConfig" withExtension:@"json"];

    NSError *error;
    [cm setApiKey:TEST_API_KEY error:&error];
    if (error) {
        return nil;
    }

    [cm setDevelopmentConfigUrl:url.absoluteString];
    error = [cm startReturningError];
    if (error) {
        NSLog(@"error starting CM: %@", error);
        return nil;
    }

    return cm;
}

- (void)testNamedCondition {
    CriticalMoments *cm = [self buildAndStartCMForTest];
    XCTAssert(cm, @"Startup issue");

    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    XCTestExpectation *expectation1 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation1];
    [cm checkNamedCondition:@"trueCondition"
                    handler:^(bool result, NSError *error) {
                      if (result && !error) {
                          [expectation1 fulfill];
                      }
                    }];

    XCTestExpectation *expectation2 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation2];
    [cm checkNamedCondition:@"falseCondition"
                    handler:^(bool result, NSError *_Nullable error) {
                      if (!result && !error) {
                          [expectation2 fulfill];
                      }
                    }];

    // check errors return errors
    XCTestExpectation *expectation3 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation3];
    [cm checkNamedCondition:@"invalidCondition"
                    handler:^(bool result, NSError *_Nullable error) {
                      if (!result && error) {
                          [expectation3 fulfill];
                      }
                    }];

    [self waitForExpectations:expectations timeout:5.0];
}

- (void)testDefaultTheme {
    CriticalMoments *cm = [self buildAndStartCMForTest];
    XCTAssert(cm, @"Startup issue");

    // Check theme loaded from config
    XCTAssert([UIColor.redColor isEqual:cm.currentTheme.bannerBackgroundColor],
              @"Default theme should have loaded bg from config");
    XCTAssert([UIColor.greenColor isEqual:cm.currentTheme.bannerForegroundColor],
              @"Default theme should have loaded fg from config");

    // Check the global sharedInstance != this instance
    XCTAssert(cm.currentTheme != CMTheme.current, @"CM instance impacted sharedInstance");
}

- (void)testSendEventBeforeStart {
    // TODO: this previously ran on simulators and CI, should try to restore this test.
    if (getenv("SIMULATOR_MODEL_IDENTIFIER") != NULL) {
        XCTSkip("not working on CI, use HW specific tests here");
    }

    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];

    NSString *randEventName = [NSString stringWithFormat:@"event_%d", arc4random()];

    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    // Inverted means we check that we don't run before we start, and queue works
    XCTestExpectation *expectationNotRun = [[XCTestExpectation alloc] init];
    expectationNotRun.inverted = true;

    // tracks that sends event after we start
    XCTestExpectation *expectation1 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation1];

    // Check order of run is order of is in order called
    NSLock *lock = [[NSLock alloc] init];
    NSMutableArray<NSNumber *> *orderRan = [[NSMutableArray alloc] init];

    // should run async after start, and not crash
    [cm sendEvent:@"custom_event"
          builtIn:false
          handler:^(NSError *_Nullable error) {
            [lock lock];
            [orderRan addObject:@1];
            [lock unlock];
            [expectationNotRun fulfill];
            XCTAssertNil(error, @"Error sending event before start");
            [expectation1 fulfill];
          }];

    // Shouldn't run yet, even if we wait 1s
    [self waitForExpectations:@[ expectationNotRun ] timeout:1.0];

    [self startCMForTest:cm];

    // should run async and not crash
    // tracks that sends event after we start
    XCTestExpectation *expectation2 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation2];
    [cm sendEvent:randEventName
          builtIn:false
          handler:^(NSError *_Nullable error) {
            [lock lock];
            [orderRan addObject:@2];
            [lock unlock];
            XCTAssertNil(error, @"failed to send rand event");
            [expectation2 fulfill];
          }];

    // should run async and not crash
    // should error because this event name is not allowed from client (built in)
    XCTestExpectation *expectation3 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation3];
    [cm sendEvent:DatamodelAppStartBuiltInEvent
          builtIn:false
          handler:^(NSError *_Nullable error) {
            [lock lock];
            [orderRan addObject:@3];
            [lock unlock];
            XCTAssertNotNil(error, @"failed to error on reserved event name");
            [expectation3 fulfill];
          }];

    // both should process
    [self waitForExpectations:expectations timeout:5.0];
    // Should process in order
    XCTAssert(orderRan.count == 3, @"all did not run");
    XCTAssert([@1 isEqualToNumber:orderRan.firstObject], @"ran out of order");
    XCTAssert([@2 isEqualToNumber:orderRan[1]], @"ran out of order");
    XCTAssert([@3 isEqualToNumber:orderRan[2]], @"ran out of order");

    XCTestExpectation *expectCount = [[XCTestExpectation alloc] init];
    // Test a condition counting events
    NSString *testCondition = [NSString stringWithFormat:@"eventCount('%@') == 1", randEventName];
    [cm checkInternalTestCondition:testCondition
                           handler:^(bool result, NSError *_Nullable error) {
                             if (result && !error) {
                                 [expectCount fulfill];
                             }
                           }];
    [self waitForExpectations:@[ expectCount ] timeout:5.0];
}

- (void)testPerformActionBeforeStart {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];
    [cm disableUserNotifications];

    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    // Inverted means we check that we don't run before we start, and queue works
    XCTestExpectation *expectationNotRun = [[XCTestExpectation alloc] init];
    expectationNotRun.inverted = true;

    // tracks that perform works after we start
    XCTestExpectation *expectationSuccess1 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectationSuccess1];

    // should run async after start, and not crash
    [cm performNamedAction:@"reviewAction"
                   handler:^(NSError *_Nullable error) {
                     [expectationNotRun fulfill];
                     XCTAssertNil(error, @"review action error");
                     [expectationSuccess1 fulfill];
                   }];

    // Shouldn't run yet, even if we wait 1s
    [self waitForExpectations:@[ expectationNotRun ] timeout:1.0];

    [self startCMForTest:cm];

    // should run async and not crash
    XCTestExpectation *expectationSuccess2 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectationSuccess2];
    [cm performNamedAction:@"reviewAction"
                   handler:^(NSError *_Nullable error) {
                     XCTAssertNil(error, @"review action error");
                     [expectationSuccess2 fulfill];
                   }];

    // should run async, and we expect error since action name not in config
    XCTestExpectation *expectationFail3 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectationFail3];
    [cm performNamedAction:@"actionWhichDoesNotExist"
                   handler:^(NSError *_Nullable error) {
                     XCTAssertNotNil(error, @"missing action did not error");
                     [expectationFail3 fulfill];
                   }];

    // confirm all are run after we start
    [self waitForExpectations:expectations timeout:5.0];
}

- (void)testCheckConditionBeforeStart {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];

    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    // Inverted means we check that we don't run before we start, and queue works
    XCTestExpectation *expectationNotRun = [[XCTestExpectation alloc] init];
    expectationNotRun.inverted = true;

    // tracks that condition works after start
    XCTestExpectation *expectationSuccess1 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectationSuccess1];

    // should run async after start, and not crash
    [cm checkInternalTestCondition:@"true"
                           handler:^(bool result, NSError *_Nullable error) {
                             [expectationNotRun fulfill];
                             if (result && !error) {
                                 [expectationSuccess1 fulfill];
                             }
                           }];

    // Shouldn't run yet, even if we wait 1s
    [self waitForExpectations:@[ expectationNotRun ] timeout:1.0];

    [self startCMForTest:cm];

    XCTestExpectation *expectationSuccess2 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectationSuccess2];

    // should run async and not crash
    [cm checkInternalTestCondition:@"false"
                           handler:^(bool result, NSError *_Nullable error) {
                             if (!error && !result) {
                                 [expectationSuccess2 fulfill];
                             }
                           }];

    // Both should have run, and returned correct results
    [self waitForExpectations:expectations timeout:5.0];
}

- (void)testHardcodedEnumCompatibility {
    // Allow compiling on iOS 16 SDK, but really should use iOS 17
#if __IPHONE_OS_VERSION_MAX_ALLOWED >= 170000
    if (@available(iOS 17.0, *)) {
        if (EKAuthorizationStatusAuthorized != EKAuthorizationStatusFullAccess) {
            XCTAssert(
                false,
                @"Code assumes EKAuthorizationStatusAuthorized == EKAuthorizationStatusFullAccess for SDK back compat");
        }

        if (4 != EKAuthorizationStatusWriteOnly) {
            XCTAssert(false, @"Code assumes EKAuthorizationStatusWriteOnly == 4 for SDK back compat");
        }
    }

    if (@available(iOS 18.0, *)) {
        if (CNAuthorizationStatusLimited != 4) {
            XCTAssert(false, @"Code assumes CNAuthorizationStatusLimited == 4");
        }
    }
#endif
}

- (void)testRegisteringProperties {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];

    // Registering built in properties should fail
    NSError *error;
    BOOL success = [cm setStringProperty:@"hello" forKey:@"platform" error:&error];
    XCTAssert(!success && error != nil, @"did not error on built in property");
    error = nil;

    // Register well known property with wrong type should fail
    success = [cm setStringProperty:@"hello" forKey:@"user_signup_date" error:&error];
    XCTAssert(!success && error, @"did not error on type missmatch");
    error = nil;

    // Register well known property with correct type should work, each type
    NSDate *signupDate = [NSDate dateWithTimeIntervalSince1970:1698093984];
    success = [cm setTimeProperty:signupDate forKey:@"user_signup_date" error:&error];
    XCTAssert(success && !error, @"failed to register well known property");
    success = [cm setStringProperty:@"test_source" forKey:@"referral_source" error:&error];
    XCTAssert(success && !error, @"failed to register well known property");
    success = [cm setBoolProperty:true forKey:@"user_signed_in" error:&error];
    XCTAssert(success && !error, @"failed to register well known property");
    success = [cm setIntegerProperty:42 forKey:@"user_referral_count" error:&error];
    XCTAssert(success && !error, @"failed to register well known property");
    success = [cm setFloatProperty:3.14 forKey:@"total_purchase_value" error:&error];
    XCTAssert(success && !error, @"failed to register well known property");

    // Registering custom propety should work
    success = [cm setStringProperty:@"hello" forKey:@"stringy" error:&error];
    XCTAssert(success && !error, @"failed to register custom property");

    NSString *jsonString = @"{\"js\": \"a\", \"jb\": true, \"jn\": 3.3}";
    success = [cm setPropertiesFromJson:[jsonString dataUsingEncoding:NSUTF8StringEncoding] error:&error];
    XCTAssert(success && !error, @"failed to register json properties");

    [self startCMForTest:cm];

    // Fetching set properties should work (both short and long form accessors)
    XCTestExpectation *expectation = [[XCTestExpectation alloc] init];
    [cm checkInternalTestCondition:
            @"user_signup_date == unixTimeSeconds(1698093984) && stringy =='hello' && custom_stringy "
            @"== 'hello' && "
            @"set_after == nil && js == 'a' && jb == true && jn == 3.3 && referral_source == 'test_source' && "
            @"user_signed_in && user_referral_count == 42 && total_purchase_value == 3.14"
                           handler:^(bool result, NSError *_Nullable er2) {
                             XCTAssert(!er2, @"test condition errored");
                             XCTAssert(result, @"test condition false");
                             [expectation fulfill];
                           }];
    // Wait as we have a race with setStringProperty below
    [self waitForExpectations:@[ expectation ] timeout:5.0];

    // registering after start should work, and quering it should return updated value
    success = [cm setStringProperty:@"hello" forKey:@"set_after" error:&error];
    XCTAssert(success && !error, @"allowed registration after start");

    XCTestExpectation *waitRegisterAfter = [[XCTestExpectation alloc] init];
    [cm checkInternalTestCondition:@"set_after == 'hello'"
                           handler:^(bool result, NSError *_Nullable er2) {
                             XCTAssert(!er2, @"test condition errored");
                             XCTAssert(result, @"test condition false");
                             [waitRegisterAfter fulfill];
                           }];

    // Both should have run, and returned correct results
    [self waitForExpectations:@[ waitRegisterAfter ] timeout:5.0];
}

- (void)testRegisteringPropertiesLegacy {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];
    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    // Registering built in properties should fail
    NSError *error;
    [cm registerStringProperty:@"hello" forKey:@"platform" error:&error];
    XCTAssertNotNil(error, @"did not error on built in property");
    error = nil;

    // Register well known property with wrong type should fail
    [cm registerStringProperty:@"hello" forKey:@"user_signup_date" error:&error];
    XCTAssertNotNil(error, @"did not error on type missmatch");
    error = nil;

    // Register well known property with correct type should work, each type
    NSDate *signupDate = [NSDate dateWithTimeIntervalSince1970:1698093984];
    [cm registerTimeProperty:signupDate forKey:@"user_signup_date" error:&error];
    XCTAssertNil(error, @"failed to register well known property");
    [cm registerStringProperty:@"test_source" forKey:@"referral_source" error:&error];
    XCTAssertNil(error, @"failed to register well known property");
    [cm registerBoolProperty:true forKey:@"user_signed_in" error:&error];
    XCTAssertNil(error, @"failed to register well known property");
    [cm registerIntegerProperty:42 forKey:@"user_referral_count" error:&error];
    XCTAssertNil(error, @"failed to register well known property");
    [cm registerFloatProperty:3.14 forKey:@"total_purchase_value" error:&error];
    XCTAssertNil(error, @"failed to register well known property");

    // Registering custom propety should work
    [cm registerStringProperty:@"hello" forKey:@"stringy" error:&error];
    XCTAssertNil(error, @"failed to register custom property");

    NSString *jsonString = @"{\"js\": \"a\", \"jb\": true, \"jn\": 3.3}";
    [cm registerPropertiesFromJson:[jsonString dataUsingEncoding:NSUTF8StringEncoding] error:&error];
    XCTAssertNil(error, @"failed to register json properties");

    [self startCMForTest:cm];

    // Fetching set properties should work (both short and long form accessors)
    XCTestExpectation *wait = [[XCTestExpectation alloc] init];
    [expectations addObject:wait];
    [cm checkInternalTestCondition:
            @"user_signup_date == unixTimeSeconds(1698093984) && stringy =='hello' && custom_stringy "
            @"== 'hello' && "
            @"stringy2 == nil && js == 'a' && jb == true && jn == 3.3 && referral_source == 'test_source' && "
            @"user_signed_in && user_referral_count == 42 && total_purchase_value == 3.14"
                           handler:^(bool result, NSError *_Nullable er2) {
                             XCTAssert(!er2, @"test condition errored");
                             XCTAssert(result, @"test condition false");
                             [wait fulfill];
                           }];

    // Both should have run, and returned correct results
    [self waitForExpectations:expectations timeout:5.0];
}

- (void)testTimezoneOffset {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];
    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    [self startCMForTest:cm];

    XCTestExpectation *expectation = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation];

    NSTimeZone *tz = NSTimeZone.localTimeZone;
    NSString *condition = [NSString stringWithFormat:@"timezone_gmt_offset == %ld", tz.secondsFromGMT];
    [cm checkInternalTestCondition:condition
                           handler:^(bool result, NSError *_Nullable er2) {
                             if (er2 || !result) {
                                 XCTAssert(false, @"timezone property failed");
                             }
                             [expectation fulfill];
                           }];

    [self waitForExpectations:expectations timeout:5.0];
}

- (void)testWeatherProviderCases:(NSArray<NSString *> *)weatherTests {
    // TODO: this doesn't work reliably in CI. Leave it as a HW only test.
    if (getenv("SIMULATOR_MODEL_IDENTIFIER") != NULL) {
        XCTSkip("Not reliably working on CI. HW only test");
    }

    CriticalMoments *cm = [self buildAndStartCMForTest];
    // TODO P2: global in test not ideal
    [CriticalMoments.sharedInstance setApiKey:TEST_API_KEY error:nil];
    XCTAssert(cm, @"Startup issue");
    @try {
        // Toronto
        [CMWeatherPropertyProvider setTestLocationOverride:[[CLLocation alloc] initWithLatitude:43.651070
                                                                                      longitude:-79.347015]];

        NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

        for (NSString *condition in weatherTests) {

            XCTestExpectation *expectation = [[XCTestExpectation alloc] init];
            [expectations addObject:expectation];
            [cm checkInternalTestCondition:condition
                                   handler:^(bool result, NSError *error) {
                                     if (!result || error) {
                                         XCTAssert(false, @"Weather test failed: %@", condition);
                                     }
                                     [expectation fulfill];
                                   }];
        }

        [self waitForExpectations:expectations timeout:15.0];
    } @finally {
        [CMWeatherPropertyProvider setTestLocationOverride:nil];
    }
}

- (void)testWeatherProviderAccurate {
    NSArray<NSString *> *weatherTests = @[
        @"weather_temperature >= -40.0 && weather_temperature <= 50.0",
        @"weather_apparent_temperature >= -40.0 && weather_apparent_temperature <= 50.0", // add_test_count
        @"weather_condition != nil && len(weather_condition) > 0",                        // add_test_count
        @"weather_cloud_cover >= 0.0 && weather_cloud_cover <= 1.0",                      // add_test_count
        @"is_daylight in ['unknown', 'daylight', 'not_daylight']",                        // add_test_count
    ];

    [self testWeatherProviderCases:weatherTests];
}

- (void)testWeatherProviderApproximate {
    NSArray<NSString *> *weatherTests = @[
        @"weather_approx_location_temperature >= -40.0 && weather_approx_location_temperature <= 50.0",
        @"weather_approx_location_apparent_temperature >= -40.0 && weather_approx_location_apparent_temperature <= 50.0", // add_test_count
        @"weather_approx_location_condition != nil && len(weather_approx_location_condition) > 0",   // add_test_count
        @"weather_approx_location_cloud_cover >= 0.0 && weather_approx_location_cloud_cover <= 1.0", // add_test_count
        @"approx_location_is_daylight in ['unknown', 'daylight', 'not_daylight']",                   // add_test_count
    ];

    [self testWeatherProviderCases:weatherTests];
}

- (void)testSharedInstance {
    CriticalMoments *cm = [CriticalMoments sharedInstance];
    XCTAssertNotNil(cm, @"shared instance not created");
    CriticalMoments *cmShared = [CriticalMoments shared];
    XCTAssertEqual(cm, cmShared, @"shared instance mismatch");
    return cm;
}

- (void)testConfigEscaping {
    // https://github.com/CriticalMoments/CriticalMoments/issues/110
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];
    [cm disableUserNotifications];
    NSBundle *testBundle = [NSBundle bundleForClass:self.class];
    NSURL *resourceBundleId =
        [testBundle.bundleURL URLByAppendingPathComponent:@"CriticalMoments_CriticalMomentsTests.bundle"];
    NSBundle *resourceBundle = [NSBundle bundleWithURL:resourceBundleId];

    BOOL success = [cm setDevelopmentConfigNameWithSuccess:@"TestResources/testConfig.json" fromBundle:resourceBundle];
    XCTAssert(success, @"Existing config file should be found");
    success = [cm setDevelopmentConfigNameWithSuccess:@"TestResources/fakeConfig.json" fromBundle:resourceBundle];
    XCTAssert(!success, @"Non existent config file should be not found");
    success = [cm setDevelopmentConfigNameWithSuccess:@"TestResources/testConfig withspace.json"
                                           fromBundle:resourceBundle];
    XCTAssert(success, @"Existing config file with charaters to escape should be found");
}

@end
