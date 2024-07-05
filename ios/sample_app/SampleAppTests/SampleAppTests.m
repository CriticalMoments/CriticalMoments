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

    // Check the furute notification is scheduled
    XCTAssert([self notificationScheduled:@"io.criticalmoments.notifications.futureNotification"],
              @"future notification not scheduled");
    // Check the past notification is not
    XCTAssert(![self notificationScheduled:@"io.criticalmoments.notifications.pastDueNotification"],
              @"past notification scheduled");
}

- (void)testNotifcationEventsAndCleanup {
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

    // Manually add a notification, as if it had been leftover from a prior config
    [self scheduleNotificationWithId:@"io.criticalmoments.notifications.priorNotification"];
    // Add a notification as if the app added it, we shouldn't mess with this
    [self scheduleNotificationWithId:@"app_notification"];

    // Check the manual notification is scheduled
    XCTAssert([self notificationScheduled:@"io.criticalmoments.notifications.priorNotification"],
              @"prior notification not scheduled");
    // Check the app notification is scheduled
    XCTAssert([self notificationScheduled:@"app_notification"], @"app notification not scheduled");
    // Check the event triggered is not scheduled as the event hasn't fired yet
    XCTAssert(![self notificationScheduled:@"io.criticalmoments.notifications.eventTriggeredNotification"],
              @"event notification scheduled too soon");

    // Send an event that triggers a notificaiton refresh
    [CriticalMoments.sharedInstance sendEvent:@"trigger_notificaiton_event"];

    // Wait for propigation
    sleep(2.0);

    // Check we unscheduled the prior notification in our cleanup
    XCTAssert(![self notificationScheduled:@"io.criticalmoments.notifications.priorNotification"],
              @"prior notification should be unscheduled");
    // Check the app notification is still scheduled
    XCTAssert([self notificationScheduled:@"app_notification"], @"app notification not scheduled");
    // Check the event trigger worked
    XCTAssert([self notificationScheduled:@"io.criticalmoments.notifications.eventTriggeredNotification"],
              @"event notif not scheduled");
    UNCalendarNotificationTrigger *initialEventRequest =
        (UNCalendarNotificationTrigger *)[self
            notificationScheduled:@"io.criticalmoments.notifications.eventTriggeredNotification"]
            .trigger;
    NSDate *initialTime = [initialEventRequest nextTriggerDate];

    // Send notification that pushes back time of the notification
    [CriticalMoments.sharedInstance sendEvent:@"trigger_notificaiton_event"];
    // Wait for propigation
    sleep(2.0);

    UNCalendarNotificationTrigger *laterEventRequest =
        (UNCalendarNotificationTrigger *)[self
            notificationScheduled:@"io.criticalmoments.notifications.eventTriggeredNotification"]
            .trigger;
    NSDate *latestTime = [laterEventRequest nextTriggerDate];
    XCTAssert(initialEventRequest != laterEventRequest, @"event notification not updated");
    XCTAssert(initialTime != nil && latestTime != nil, @"date issue");
    XCTAssert([initialTime compare:latestTime] == NSOrderedAscending, @"event did not push back delviery time");
    NSTimeInterval timediff = [latestTime timeIntervalSinceDate:initialTime];
    XCTAssert(timediff > 1.8 && timediff < 2.2, @"event did not push back delivery time by correct amount (2s)");

    // Send an event that triggers cancelation of notification
    [CriticalMoments.sharedInstance sendEvent:@"cancel_notification"];

    // Wait for propigation
    sleep(2.0);

    // Check we unscheduled the prior notification in our cleanup
    XCTAssert(![self notificationScheduled:@"io.criticalmoments.notifications.priorNotification"],
              @"prior notification should be unscheduled");
    // Check the app notification is still scheduled
    XCTAssert([self notificationScheduled:@"app_notification"], @"app notification not scheduled");
    // Check the event trigger worked
    XCTAssert(![self notificationScheduled:@"io.criticalmoments.notifications.eventTriggeredNotification"],
              @"event notif is scheduled after cancelation");
}

- (void)scheduleNotificationWithId:(NSString *)notifId {
    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
    content.title = @"Test Notification";
    UNTimeIntervalNotificationTrigger *trigger = [UNTimeIntervalNotificationTrigger triggerWithTimeInterval:600
                                                                                                    repeats:NO];
    UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:notifId
                                                                          content:content
                                                                          trigger:trigger];
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center addNotificationRequest:request withCompletionHandler:nil];
}

- (UNNotificationRequest *)notificationScheduled:(NSString *)notifId {
    UNNotificationRequest *__block notifRequest;
    XCTestExpectation *expectation = [[XCTestExpectation alloc] init];

    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center
        getPendingNotificationRequestsWithCompletionHandler:^(NSArray<UNNotificationRequest *> *_Nullable requests) {
          for (UNNotificationRequest *request in requests) {

              if ([request.identifier isEqualToString:notifId]) {
                  if (notifRequest != nil) {
                      XCTAssert(false, @"Two notifications with same id");
                  }
                  notifRequest = request;
              }
          }
          [expectation fulfill];
        }];

    [self waitForExpectations:@[ expectation ] timeout:2.0];
    return notifRequest;
}

@end
