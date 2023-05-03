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

- (void)testIntegration {
    NSString *pongResponse = [CriticalMoments objcPing];
    XCTAssert([@"objcPong" isEqual:pongResponse], @"CM integration broken");

    NSString *goPongResponse = [CriticalMoments goPing];
    XCTAssert([@"AppcorePong->PongCmCore" isEqual:goPongResponse],
              @"CM Go integration broken");
}

- (void)testPerformanceExample {
    // This is an example of a performance test case.
    [self measureBlock:^{
        // Put the code you want to measure the time of here.
    }];
}

@end
