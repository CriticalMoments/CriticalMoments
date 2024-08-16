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
    XCTAssert([black isEqual:[CMUtils colorFromHexString:@"#000000"]], @"not black enough");
    UIColor *white = [[UIColor alloc] initWithRed:1 green:1 blue:1 alpha:1];
    XCTAssert([white isEqual:[CMUtils colorFromHexString:@"#ffffff"]], @"not white enough");
    XCTAssert([UIColor.blueColor isEqual:[CMUtils colorFromHexString:@"#0000ff"]], @"not blue enough");
    XCTAssert([UIColor.greenColor isEqual:[CMUtils colorFromHexString:@"#00ff00"]], @"not green enough");
    XCTAssert([UIColor.redColor isEqual:[CMUtils colorFromHexString:@"#ff0000"]], @"not red enough");
}

- (void)testLocationNoising {
    int totalLocations = 500;
    int locationsPerRun = 300;
    CLLocationDistance totalDistanceOffset = 0.0;
    CLLocationDistance totalAvgOffsets = 0.0;
    int amountOfNoise = 1000; // 1km

    for (int i = 0; i < totalLocations; i++) {
        float randLat = (((float)arc4random() / UINT32_MAX) * 180.0) - 90.0;
        XCTAssert(randLat >= -90.0);
        XCTAssert(randLat <= 90.0);
        float randLong = (((float)arc4random() / UINT32_MAX) * 230.0) - 180.0;
        XCTAssert(randLong >= -180.0);
        XCTAssert(randLong <= 180.0);
        CLLocation *start = [[CLLocation alloc] initWithLatitude:randLat longitude:randLong];

        double totalLat = 0.0;
        double totalLong = 0.0;

        for (int ii = 0; ii < locationsPerRun; ii++) {

            CLLocation *end = [CMUtils noiseLocation:start maxNoise:amountOfNoise];

            // check distance within amountOfNoise (with floating point error buffer)
            CLLocationDistance dist = [start distanceFromLocation:end];
            if (dist > amountOfNoise * 1.05) {
                NSLog(@"dist: %f", dist);
                XCTAssert(false, @"distance too far");
            }

            // save values for final check
            totalDistanceOffset += dist;
            totalLat += end.coordinate.latitude;
            totalLong += end.coordinate.longitude;
        }

        // The avg center point of all the samples should be close(ish) to start, as noise should average out.
        // This ensures we're not always going in the same fixed direction, which would be reversible.
        double avgLat = totalLat / (double)(locationsPerRun);
        double avgLong = totalLong / (double)(locationsPerRun);
        CLLocation *avgNoisedLocation = [[CLLocation alloc] initWithLatitude:avgLat longitude:avgLong];
        CLLocationDistance dist = [start distanceFromLocation:avgNoisedLocation];
        totalAvgOffsets += dist;
    }

    // Confirm avg distance offset is roughly half of noise param (with fudge factor)
    // This should be true because we add linear noise
    double avgDistOffset = totalDistanceOffset / (locationsPerRun * totalLocations);
    double diffFromExpected = (amountOfNoise / 2) - avgDistOffset;
    // This is probalistic, but with n=150k shouldn't really fail. If it fails 1 our of a thousands times don't sweat
    // it.
    XCTAssert(diffFromExpected < (0.05 * amountOfNoise) && diffFromExpected > (-0.05 * amountOfNoise));

    // Confirm that on average, the center point of there noised samples is close to the point we're starting with
    // This is probalistic, but with n=150k shouldn't really fail. If it fails 1 our of a thousands times don't sweat
    // it.
    XCTAssert(totalAvgOffsets / (double)totalLocations < amountOfNoise * 0.1);
}

@end
