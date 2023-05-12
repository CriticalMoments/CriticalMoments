//
//  SnapshotTests.m
//  SampleAppTests
//
//  Created by Steve Cosman on 2023-05-11.
//

#import <XCTest/XCTest.h>

@import iOSSnapshotTestCase;
@import iOSSnapshotTestCaseCore;

#import "SampleAppCoreViewController.h"
#import "Utils.h"

@interface SnapshotTests : FBSnapshotTestCase

@end

@implementation SnapshotTests

- (void)setUp {
    [super setUp];

    // record new screenshots
    // self.recordMode = YES;

    // Needed for system UI and transparency
    self.usesDrawViewHierarchyInRect = YES;
}

- (void)tearDown {
    [super tearDown];
}

- (void)testScreenshotAllSampleAppFeatures {
    UIWindow *window = [[UIWindow alloc] init];
    SampleAppCoreViewController *mainVc =
        [[SampleAppCoreViewController alloc] init];
    window.rootViewController = mainVc;
    [window makeKeyAndVisible];

    [[NSRunLoop currentRunLoop]
        runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];

    FBSnapshotVerifyView([Utils keyWindow], @"startup");

    // make animations super fast (20x)
    [Utils keyWindow].layer.speed = 20.0;

    [self recursiveActionPlayer:mainVc.demoRoot];
}

- (void)recursiveActionPlayer:(CMDemoScreen *)demoscreen {
    // "Play" each action in each section, screenshot, revert changes,
    // then move on. Goes into menus recursively when it finds one.

    // Hide the scroll bar as it can cause tests to fail
    UINavigationController *navController = [Utils appNavControl];
    UIViewController *vc = navController.visibleViewController;
    if ([vc isKindOfClass:UITableViewController.class]) {
        UITableViewController *tableVc = (UITableViewController *)vc;
        tableVc.tableView.showsHorizontalScrollIndicator = NO;
        tableVc.tableView.showsVerticalScrollIndicator = NO;
    }

    for (CMDemoSection *section in demoscreen.sections) {
        for (CMDemoAction *action in section.actions) {
            if (action.skipInUiTesting) {
                continue;
            }

            // Perform the action and give it time to render
            [action performAction];
            [[NSRunLoop currentRunLoop]
                runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.4]];

            if (action.actionNextScreen) {
                // this is a nav action, follow it!
                [self recursiveActionPlayer:action.actionNextScreen];

                // Pop back to this VC
                UINavigationController *navController = [Utils appNavControl];
                [navController popViewControllerAnimated:NO];
            } else {
                // Snapshot test!
                FBSnapshotVerifyView([Utils keyWindow], action.title);
            }

            // reset state for next test, and give it time to render
            [action resetForTests];
            [[NSRunLoop currentRunLoop]
                runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.2]];
        }
    }
}

@end
