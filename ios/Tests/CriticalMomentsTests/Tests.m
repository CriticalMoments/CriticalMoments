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
    [self startCMForTest:cm];
    return cm;
}

- (void)startCMForTest:(CriticalMoments *)cm {
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
    CriticalMoments *cm = [self buildAndStartCMForTest];
    XCTAssert(cm, @"Startup issue");

    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    XCTestExpectation *expectation1 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation1];
    [cm checkNamedCondition:@"trueCondition"
                  condition:@"true"
                    handler:^(bool result, NSError *error) {
                      if (result && !error) {
                          [expectation1 fulfill];
                      }
                    }];

    XCTestExpectation *expectation2 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation2];
    [cm checkNamedCondition:@"falseCondition"
                  condition:@"false"
                    handler:^(bool result, NSError *_Nullable error) {
                      if (!result && !error) {
                          [expectation2 fulfill];
                      }
                    }];

    // check errors return errors
    XCTestExpectation *expectation3 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation3];
    [cm checkNamedCondition:@"invalidCondition"
                  condition:@"fake_var > 6"
                    handler:^(bool result, NSError *_Nullable error) {
                      if (!result && error) {
                          [expectation3 fulfill];
                      }
                    }];

    // Override this by name the test json config file, should return true.
    XCTestExpectation *expectation4 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation4];
    [cm checkNamedCondition:@"overrideToTrue"
                  condition:@"false"
                    handler:^(bool result, NSError *_Nullable error) {
                      if (result && !error) {
                          [expectation4 fulfill];
                      }
                    }];

    // This should show a warning in debug mode, but should pass
    // "falseCondition" name conflict with early use and different condition
    XCTestExpectation *expectation5 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectation5];
    [cm checkNamedCondition:@"falseCondition"
                  condition:@"true"
                    handler:^(bool result, NSError *_Nullable error) {
                      if (result && !error) {
                          [expectation5 fulfill];
                      }
                    }];

    [self waitForExpectations:expectations timeout:5.0];
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

- (void)testSendEventBeforeStart {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];

    // should run async after start, and not crash
    [cm sendEvent:@"eventName"];

    [NSThread sleepForTimeInterval:5.0];

    [self startCMForTest:cm];

    // should run async and not crash
    [cm sendEvent:@"eventName2"];

    // TODO both should process, in right order
}

- (void)testPerformActionBeforeStart {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];

    // should run async after start, and not crash
    NSError *error;
    [cm performNamedAction:@"reviewAction" error:&error];
    // TODO: fails now, but should work once we make this API async
    XCTAssert(error, @"perform action before start errored");

    [self startCMForTest:cm];

    // should run async and not crash
    error = nil;
    [cm performNamedAction:@"reviewAction" error:&error];
    XCTAssert(!error, @"condition after start errored");

    // TODO both should process, in right order
}

- (void)testCheckConditionBeforeStart {
    CriticalMoments *cm = [[CriticalMoments alloc] initInternal];

    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    // Inverted means we check that we don't run before we start, and queue works
    XCTestExpectation *expectationRun = [[XCTestExpectation alloc] init];
    expectationRun.inverted = true;

    // tracks that condition works after start
    XCTestExpectation *expectationSuccess1 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectationSuccess1];

    // should run async after start, and not crash
    [cm checkNamedCondition:@"nonName"
                  condition:@"true"
                    handler:^(bool result, NSError *_Nullable error) {
                      [expectationRun fulfill];
                      if (result && !error) {
                          [expectationSuccess1 fulfill];
                      }
                    }];

    // Shouldn't run yet, even if we wait 1s
    [self waitForExpectations:@[ expectationRun ] timeout:1.0];

    [self startCMForTest:cm];

    XCTestExpectation *expectationSuccess2 = [[XCTestExpectation alloc] init];
    [expectations addObject:expectationSuccess2];

    // should run async and not crash
    [cm checkNamedCondition:@"nonName2"
                  condition:@"false"
                    handler:^(bool result, NSError *_Nullable error) {
                      if (!error && !result) {
                          [expectationSuccess2 fulfill];
                      }
                    }];

    // Both should have run, and returned correct results
    [self waitForExpectations:expectations timeout:5.0];
}
@end
