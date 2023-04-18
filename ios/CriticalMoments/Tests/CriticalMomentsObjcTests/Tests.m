//
//  CriticalMomentsTests.m
//  CriticalMomentsTests
//
//  Created by scosman on 04/12/2023.
//  Copyright (c) 2023 scosman. All rights reserved.
//

@import XCTest;

@import CriticalMomentsObjc;
#import "CriticalMoments.h"

@interface Tests : XCTestCase

@end

@implementation Tests

- (void)setUp
{
    [super setUp];
    // Put setup code here. This method is called before the invocation of each test method in the class.
}

- (void)tearDown
{
    // Put teardown code here. This method is called after the invocation of each test method in the class.
    [super tearDown];
}

- (void)testPing
{
    NSString *response = [CriticalMoments objcPing];
    XCTAssert([@"objcPong" isEqualToString:response], @"Expected ping to pong -- objective C tests not working end to end");
}

@end

