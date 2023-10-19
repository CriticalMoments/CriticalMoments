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

    // Wait for main thread to start responding. Needed for CI or we hit timeout below.
    dispatch_semaphore_t mainWait = dispatch_semaphore_create(0);
    dispatch_async(dispatch_get_main_queue(), ^{
      dispatch_semaphore_signal(mainWait);
    });
    [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.5]];
    dispatch_semaphore_wait(mainWait, dispatch_time(DISPATCH_TIME_NOW, 120.0 * NSEC_PER_SEC));

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
    [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.5]];

    [self waitForExpectations:expectations timeout:20.0];
}

@end
