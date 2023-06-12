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
}

- (void)tearDown {
    // Put teardown code here. This method is called after the invocation of
    // each test method in the class.
}

- (void)testBasicIntegration {
    NSString *pongResponse = [CriticalMoments objcPing];
    XCTAssert([@"objcPong" isEqualToString:pongResponse], @"CM integration broken");

    NSString *goPongResponse = [CriticalMoments goPing];
    XCTAssert([@"AppcorePong->PongCmCore" isEqualToString:goPongResponse], @"CM Go integration broken");
}

- (void)testDefaultTheme {
    // TOOD: this test isn't robust. Race condition with app config,
    // and edits global state impacting other tests/app.

    // Ugly wait to wait for startup of CM which is async
    XCTestExpectation *expectation = [self expectationWithDescription:@"CM startup done"];
    dispatch_async(dispatch_get_main_queue(), ^{
      // Ensure a default theme from config is loaded into app
      NSBundle *testBundle = [NSBundle bundleForClass:self.class];
      NSURL *url = [testBundle URLForResource:@"defaultThemeTest" withExtension:@"json"];
      [CriticalMoments setConfigUrl:url.absoluteString];
      [CriticalMoments start];

      dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 0.1 * NSEC_PER_SEC), dispatch_get_main_queue(), ^{
        [expectation fulfill];
      });
    });

    [self waitForExpectationsWithTimeout:10.0
                                 handler:^(NSError *error) {
                                   XCTAssert(error == nil, @"Error with test %@", error);
                                   XCTAssert([UIColor.redColor isEqual:CMTheme.current.bannerBackgroundColor],
                                             @"Default theme should have loaded bg from config");
                                   XCTAssert([UIColor.greenColor isEqual:CMTheme.current.bannerForegroundColor],
                                             @"Default theme should have loaded fg from config");
                                 }];
}

@end
