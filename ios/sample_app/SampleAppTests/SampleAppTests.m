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
}

- (void)tearDown {
}

- (void)testBasicIntegration {

    NSString *pongResponse = [CriticalMoments.sharedInstance objcPing];
    XCTAssert([@"objcPong" isEqualToString:pongResponse], @"CM integration broken");

    NSString *goPongResponse = [CriticalMoments.sharedInstance goPing];
    XCTAssert([@"AppcorePong->PongCmCore" isEqualToString:goPongResponse], @"CM Go integration broken");
}

- (void)testCanOpenUrlEndToEnd {
    CriticalMoments *cm = CriticalMoments.sharedInstance;

    NSMutableArray<XCTestExpectation *> *expectations = [[NSMutableArray alloc] init];

    NSDictionary *cases = @{
        @"testCanOpenOwnUrlScheme" : @"canOpenUrl('critical-moments-sampleapp://home') == true",
        @"testCanOpenHttpUrl" : @"canOpenUrl('http://criticalmoments.io') == true",
        @"testCantOpenInvalidUrl" : @"canOpenUrl('not a url') == false",
        @"testCantOpenUnknownScheme" : @"canOpenUrl('asfsdfdsfsdf://asdf.com') == false",
    };

    for (NSString *name in cases.keyEnumerator) {
        NSString *condition = cases[name];

        XCTestExpectation *expectation = [[XCTestExpectation alloc] initWithDescription:name];
        [expectations addObject:expectation];
        [cm checkNamedCondition:name
                      condition:condition
                        handler:^(bool result, NSError *_Nullable error) {
                          if (result && !error) {
                              [expectation fulfill];
                          }
                        }];
    }

    [self waitForExpectations:expectations timeout:1.0];
}

@end
