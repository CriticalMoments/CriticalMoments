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
    NSString *response = [CriticalMoments objcPing];
    XCTAssert(
        [@"objcPong" isEqualToString:response],
        @"Expected ping to pong -- objective C tests not working end to end");
}

- (void)testAppcoreIntegration {
    NSString *response = AppcoreGoPing();
    XCTAssert([@"AppcorePong->PongCmCore" isEqualToString:response],
              @"Expected ping to pong -- Appcore framework integration not "
              @"working end to end");

    NSString *fullyIntegratedRespons = [CriticalMoments goPing];
    XCTAssert(
        [@"AppcorePong->PongCmCore" isEqualToString:fullyIntegratedRespons],
        @"Expected ping to pong -- Appcore e2e framework integration not "
        @"working end to end");
}

@end
