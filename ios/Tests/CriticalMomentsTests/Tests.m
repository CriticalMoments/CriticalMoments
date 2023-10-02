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
    NSString *response = [CriticalMoments.sharedInstance objcPing];
    XCTAssert([@"objcPong" isEqualToString:response],
              @"Expected ping to pong -- objective C tests not working end to end");
}

- (void)testAppcoreIntegration {
    NSString *response = AppcoreGoPing();
    XCTAssert([@"AppcorePong->PongCmCore" isEqualToString:response],
              @"Expected ping to pong -- Appcore framework integration not "
              @"working end to end");

    NSString *fullyIntegratedRespons = [CriticalMoments.sharedInstance goPing];
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

    [cm setConfigUrl:url.absoluteString];
    error = [cm startReturningError];
    if (error) {
        return nil;
    }

    return cm;
}

- (void)testNamedCondition {
    NSError *error;
    CriticalMoments *cm = [self buildAndStartCMForTest];
    XCTAssert(cm, @"Startup issue");

    bool result = [cm checkNamedCondition:@"trueCondition" condition:@"true" error:&error];
    XCTAssert(result, @"result should be true");
    XCTAssert(error == nil, @"error should be nil");

    error = nil;
    result = [cm checkNamedCondition:@"falseCondition" condition:@"false" error:&error];
    XCTAssert(!result, @"result should be false");
    XCTAssert(error == nil, @"error should be nil");

    error = nil;
    result = [cm checkNamedCondition:@"invalidCondition" condition:@"fake_var > 6" error:&error];
    XCTAssert(!result, @"result should be false");
    XCTAssert(error != nil, @"error should not be nil");

    // Override this by name the test json config file, should return true.
    error = nil;
    result = [cm checkNamedCondition:@"overrideToTrue" condition:@"false" error:&error];
    XCTAssert(result, @"result should be true");
    XCTAssert(error == nil, @"error should be nil");

    // This should show a warning in debug mode, but should pass
    // "falseCondition" name conflict with early use and different condition
    error = nil;
    result = [cm checkNamedCondition:@"falseCondition" condition:@"true" error:&error];
    XCTAssert(result, @"result should be true");
    XCTAssert(error == nil, @"error should be nil");
}

- (void)testDefaultTheme {
    CriticalMoments *cm = [self buildAndStartCMForTest];
    XCTAssert(cm, @"Startup issue");

    // Not ideal that starting any CM instance will set gloabal default theme,
    // but only CM.sharedInstance is exposed in public API so P2.
    XCTAssert([UIColor.redColor isEqual:CMTheme.current.bannerBackgroundColor],
              @"Default theme should have loaded bg from config");
    XCTAssert([UIColor.greenColor isEqual:CMTheme.current.bannerForegroundColor],
              @"Default theme should have loaded fg from config");
}

@end
