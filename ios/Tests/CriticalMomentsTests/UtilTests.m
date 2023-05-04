//
//  UtilTests.m
//
//
//  Created by Steve Cosman on 2023-05-03.
//

#import <XCTest/XCTest.h>

#import "utils/CMUtils.h"

@interface UtilTests : XCTestCase

@end

@implementation UtilTests

- (void)setUp {
    // Put setup code here. This method is called before the invocation of each
    // test method in the class.
}

- (void)tearDown {
    // Put teardown code here. This method is called after the invocation of
    // each test method in the class.
}

- (void)testHexColorUtil {
    // lenth/format mismatch
    XCTAssertNil([CMUtils colorFromHexString:@"#0000000"], @"too long");
    XCTAssertNil([CMUtils colorFromHexString:@"#00000"], @"too short");
    XCTAssertNil([CMUtils colorFromHexString:@""], @"too empty");

    // Invalid charaters
    XCTAssertNil([CMUtils colorFromHexString:@"#0000gg"], @"too invalid");
    XCTAssertNil([CMUtils colorFromHexString:@"#gg0000"], @"too invalid");

    // Valid colors
    UIColor *black = [[UIColor alloc] initWithRed:0 green:0 blue:0 alpha:1];
    XCTAssert([black isEqual:[CMUtils colorFromHexString:@"#000000"]],
              @"not black enough");
    UIColor *white = [[UIColor alloc] initWithRed:1 green:1 blue:1 alpha:1];
    XCTAssert([white isEqual:[CMUtils colorFromHexString:@"#ffffff"]],
              @"not white enough");
    XCTAssert(
        [UIColor.blueColor isEqual:[CMUtils colorFromHexString:@"#0000ff"]],
        @"not blue enough");
    XCTAssert(
        [UIColor.greenColor isEqual:[CMUtils colorFromHexString:@"#00ff00"]],
        @"not green enough");
    XCTAssert(
        [UIColor.redColor isEqual:[CMUtils colorFromHexString:@"#ff0000"]],
        @"not red enough");
}

@end
