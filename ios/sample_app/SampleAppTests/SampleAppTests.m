//
//  SampleAppTests.m
//  SampleAppTests
//
//  Created by Steve Cosman on 2023-04-22.
//

#import <XCTest/XCTest.h>

#import "../SampleApp/AppDelegate.h"
#import "../SampleApp/DemoScreens/BuiltInThemesDemoScreen.h"
#import "../SampleApp/Utils.h"
#import "UserNotifications/UserNotifications.h"

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
    id<UIApplicationDelegate> ad = UIApplication.sharedApplication.delegate;
    AppDelegate *aad = (AppDelegate *)ad;
    CriticalMoments *cm = [aad cmInstance];

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
        [cm checkInternalTestCondition:condition
                               handler:^(bool result, NSError *_Nullable error) {
                                 if (error != nil) {
                                     XCTAssert(false, @"CanOpenUrl test failed with error: %@", error);
                                 }
                                 XCTAssertTrue(result, @"CanOpenUrl test did pass for condition check: %@", name);
                                 [expectation fulfill];
                               }];
    }

    [self waitForExpectations:expectations timeout:20.0];
}

- (void)testThemeCount {
    NSDictionary *themeDescriptions = [BuiltInThemesDemoScreen themeDescriptions];
    int expected = [CriticalMoments.sharedInstance builtInBaseThemeCount];
    XCTAssert(themeDescriptions.count == expected, @"Expected %d themes in demo app, got %lu", expected,
              (unsigned long)themeDescriptions.count);
}

- (void)testBundleCheck {
    // Roundabout test to ensure urlAllowedForDebugLoad excludes writeable directories.
    // XCUnitTests have their own set of directories, so we save paths in the main app, and check them here
    BOOL success = [Utils verifyTestFileUrls];
    XCTAssert(success, @"A app-writeable directory passes urlAllowedForDebugLoad check");
}

- (void)testNotifcations {
    // Need to notification permissions for this test to work
    XCTestExpectation *approvalExpectation = [[XCTestExpectation alloc] init];
    BOOL __block approved = false;
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center getNotificationSettingsWithCompletionHandler:^(UNNotificationSettings *_Nonnull settings) {
      approved = settings.authorizationStatus == UNAuthorizationStatusAuthorized;
      [approvalExpectation fulfill];
    }];
    [self waitForExpectations:@[ approvalExpectation ] timeout:2.0];
    if (!approved) {
        XCTSkip(@"User notification permission not approved, test won't work");
    }

    XCTestExpectation *waitExpectation = [[XCTestExpectation alloc] init];

    // Check we schedule the future one, but not the past one.
    [center getPendingNotificationRequestsWithCompletionHandler:^(NSArray<UNNotificationRequest *> *_Nonnull requests) {
      XCTAssert(requests.count == 1, @"only one test notification should get scheduled. Got %ld", requests.count);
      UNNotificationRequest *n = requests.firstObject;
      XCTAssert([@"io.criticalmoments.notifications.futureNotification" isEqualToString:n.identifier],
                @"wrong notification scheduled. Got %@", n.identifier);
      XCTAssert([@"Future" isEqualToString:n.content.title], @"notification title mismatch");
      [waitExpectation fulfill];
    }];
    [self waitForExpectations:@[ waitExpectation ] timeout:10.0];
}
@end
