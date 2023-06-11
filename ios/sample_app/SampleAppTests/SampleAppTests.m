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
    XCTAssert([@"objcPong" isEqual:pongResponse], @"CM integration broken");

    NSString *goPongResponse = [CriticalMoments goPing];
    XCTAssert([@"AppcorePong->PongCmCore" isEqual:goPongResponse], @"CM Go integration broken");
}

- (void)testDefaultTheme {
    // Ensure a default theme from config is loaded into app
    NSBundle *testBundle = [NSBundle bundleForClass:self.class];
    NSURL *url = [testBundle URLForResource:@"defaultThemeTest" withExtension:@"json"];
    [CriticalMoments setConfigUrl:url.absoluteString];
    [CriticalMoments start];

    XCTAssert([UIColor.redColor isEqual:CMTheme.current.bannerBackgroundColor],
              @"Default theme should have loaded bg from config");
    XCTAssert([UIColor.greenColor isEqual:CMTheme.current.bannerForegroundColor],
              @"Default theme should have loaded fg from config");
}

- (void)testPerformanceExample {
    // This is an example of a performance test case.
    [self measureBlock:^{
        // Put the code you want to measure the time of here.
    }];
}

@end
