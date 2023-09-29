//
//  SampleAppTests.m
//  SampleAppTests
//
//  Created by Steve Cosman on 2023-04-22.
//

#import <XCTest/XCTest.h>

@import CriticalMoments;

@interface SampleAppTests : XCTestCase

@end

@implementation SampleAppTests

- (void)setUp {
    // Put setup code here. This method is called before the invocation of each
    // test method in the class.

    // Ugly wait to wait for startup of CM which is async
    XCTestExpectation *expectation = [self expectationWithDescription:@"CM startup done"];
    dispatch_async(dispatch_get_main_queue(), ^{
      // Ensure a default theme from config is loaded into app
      NSBundle *testBundle = [NSBundle bundleForClass:self.class];
      NSURL *url = [testBundle URLForResource:@"defaultThemeTest" withExtension:@"json"];
      [CriticalMoments setConfigUrl:url.absoluteString];
      [CriticalMoments start];

      dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 0.2 * NSEC_PER_SEC), dispatch_get_main_queue(), ^{
        [expectation fulfill];
      });
    });
}

- (void)tearDown {
    // Put teardown code here. This method is called after the invocation of
    // each test method in the class.
}

- (void)testBasicIntegration {
    [self waitForExpectationsWithTimeout:10.0
                                 handler:^(NSError *herr) {
                                   NSString *pongResponse = [CriticalMoments objcPing];
                                   XCTAssert([@"objcPong" isEqualToString:pongResponse], @"CM integration broken");

                                   NSString *goPongResponse = [CriticalMoments goPing];
                                   XCTAssert([@"AppcorePong->PongCmCore" isEqualToString:goPongResponse],
                                             @"CM Go integration broken");
                                 }];
}

- (void)testDefaultTheme {
    // TOOD: this test isn't robust. Race condition with app config,
    // and edits global state impacting other tests/app.

    [self waitForExpectationsWithTimeout:10.0
                                 handler:^(NSError *error) {
                                   XCTAssert(error == nil, @"Error with test %@", error);
                                   XCTAssert([UIColor.redColor isEqual:CMTheme.current.bannerBackgroundColor],
                                             @"Default theme should have loaded bg from config");
                                   XCTAssert([UIColor.greenColor isEqual:CMTheme.current.bannerForegroundColor],
                                             @"Default theme should have loaded fg from config");
                                 }];
}

- (void)testApiKeyValidation {
    [self waitForExpectationsWithTimeout:10.0
                                 handler:^(NSError *herr) {
                                   NSError *error;

                                   // TODO: this would be a lot better on non-global CM object
                                   [CriticalMoments setApiKey:@"" error:&error];
                                   XCTAssert(error != nil, @"Empty API key passed validation");

                                   error = nil;
                                   [CriticalMoments setApiKey:@"invalid" error:&error];
                                   XCTAssert(error != nil, @"Invalid API key passed validation");

                                   error = nil;
                                   [CriticalMoments
                                       setApiKey:@"CM1-aGVsbG86d29ybGQ=-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-"
                                                 @"MEUCIQCUfx6xlmQ0kdYkuw3SMFFI6WXrCWKWwetXBrXXG2hjAwIgWBPIMrdM1ET0Hbpn"
                                                 @"Xlnpj/f+VXtjRTqNNz9L/AOt4GY="
                                           error:&error];
                                   XCTAssert(error != nil, @"API key from another app passed validation");

                                   // This key is only valid for this sample app
                                   NSString *apiKey =
                                       @"CM1-Yjppby5jcml0aWNhbG1vbWVudHMuU2FtcGxlQXBw-MEYCIQCOd0JTuuUtgTJkDUsQH0EQMhJ+"
                                       @"kKysBBfjdxZKqgTBDAIhAMo/OGSysVA0iOscz+mKDqY8UizldA8sZj2a3/mAZIzB";
                                   error = nil;
                                   [CriticalMoments setApiKey:apiKey error:&error];
                                   XCTAssert(error == nil, @"API key failed validation");
                                 }];
}

- (void)testNamedCondition {
    [self waitForExpectationsWithTimeout:10.0
                                 handler:^(NSError *herr) {
                                   NSError *error;
                                   bool result = [CriticalMoments checkNamedCondition:@"trueCondition"
                                                                            condition:@"true"
                                                                                error:&error];
                                   XCTAssert(result, @"result should be true");
                                   XCTAssert(error == nil, @"error should be nil");

                                   error = nil;
                                   result = [CriticalMoments checkNamedCondition:@"falseCondition"
                                                                       condition:@"false"
                                                                           error:&error];
                                   XCTAssert(!result, @"result should be false");
                                   XCTAssert(error == nil, @"error should be nil");

                                   error = nil;
                                   result = [CriticalMoments checkNamedCondition:@"invalidCondition"
                                                                       condition:@"fake_var > 6"
                                                                           error:&error];
                                   XCTAssert(!result, @"result should be false");
                                   XCTAssert(error != nil, @"error should not be nil");

                                   // Override this by name the test json config file, should return true.
                                   error = nil;
                                   result = [CriticalMoments checkNamedCondition:@"overrideToTrue"
                                                                       condition:@"false"
                                                                           error:&error];
                                   XCTAssert(result, @"result should be true");
                                   XCTAssert(error == nil, @"error should be nil");

                                   // This should show a warning in debug mode, but should pass
                                   // "falseCondition" name conflict with early use and different condition
                                   error = nil;
                                   result = [CriticalMoments checkNamedCondition:@"falseCondition"
                                                                       condition:@"true"
                                                                           error:&error];
                                   XCTAssert(result, @"result should be true");
                                   XCTAssert(error == nil, @"error should be nil");
                                 }];
}

@end
