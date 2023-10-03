//
//  SnapshotTests.m
//  SampleAppTests
//
//  Created by Steve Cosman on 2023-05-11.
//

#import <XCTest/XCTest.h>

@import iOSSnapshotTestCase;
@import iOSSnapshotTestCaseCore;

#import "SampleAppTests-Swift.h"

#import "SampleAppCoreViewController.h"
#import "Utils.h"

@interface SnapshotTests : FBSnapshotTestCase

@end

//#define FB_REFERENCE_IMAGE_DIR                                                                                         \
    "/Users/scosman/Dropbox/workspace/criticalmoments/ios/sample_app/SnapshotTests/ReferenceImages"
// #define IMAGE_DIFF_DIR "/Users/scosman/Dropbox/workspace/criticalmoments/ios/sample_app/SnapshotTests/FailureDiffs"

@implementation SnapshotTests

- (void)setUp {
    [super setUp];

    // record new screenshots
    // self.recordMode = YES;

    // Needed for system UI and transparency
    self.usesDrawViewHierarchyInRect = YES;

    // Filename options
    self.fileNameOptions = FBSnapshotTestCaseFileNameIncludeOptionOS | FBSnapshotTestCaseFileNameIncludeOptionDevice |
                           FBSnapshotTestCaseFileNameIncludeOptionScreenScale |
                           FBSnapshotTestCaseFileNameIncludeOptionScreenSize;

    self.continueAfterFailure = true;
}

- (void)tearDown {
    [super tearDown];
}

- (void)testScreenshotAllSampleAppFeatures {
    [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
    SampleAppCoreViewController *mainVc;
    UIWindow *window = [Utils keyWindow];
    UIViewController *rvc = window.rootViewController;
    if ([rvc isKindOfClass:SampleAppCoreViewController.class]) {
        mainVc = (SampleAppCoreViewController *)rvc;
        mainVc.backgroundColor = [UIColor greenColor];
    } else {
        XCTAssertTrue(false, @"Could not get root vc");
    }

    // make animations super fast (20x)
    [Utils keyWindow].layer.speed = 20.0;

    [self recursiveActionPlayer:mainVc.demoRoot];
}

- (void)recursiveActionPlayer:(CMDemoScreen *)demoscreen {
    // "Play" each action in each section, screenshot, revert changes,
    // then move on. Goes into menus recursively when it finds one.

    // Hide the demo UI as we aren't testing it. Things like the scroll bars
    // can mess up tests, but also, as we add features we don't want tests to
    // fail because the table view updates
    UINavigationController *navController = [Utils appNavControl];
    UIViewController *vc = navController.visibleViewController;
    vc.view.alpha = 0.0;

    for (CMDemoSection *section in demoscreen.sections) {
        for (CMDemoAction *action in section.actions) {
            if (action.skipInUiTesting) {
                continue;
            }

            // Perform the action and give it time to render
            [action performAction];
            [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.4]];

            if (action.actionNextScreen) {
                // this is a nav action, follow it!
                [self recursiveActionPlayer:action.actionNextScreen];

                // Pop back to this VC
                UINavigationController *navController = [Utils appNavControl];
                [navController popViewControllerAnimated:NO];
            } else {
                // Snapshot test!
                // FBSnapshotVerifyView([Utils keyWindow], action.title);

                UIImage *screenshot = [self screenshotWindow:[Utils keyWindow]];
                CMSnapshotWrapper *w = [[CMSnapshotWrapper alloc] init];
                [w assertSnapshotImageOf:screenshot named:action.title];
                //[w assertSnapshotVCOf:[Utils keyWindow].rootViewController named:action.title];
                //[w assertSnapshotOf:[Utils keyWindow]];
                //[w assertSnapshotOf:[Utils keyWindow] named:action.title];
            }

            // reset state for next test, and give it time to render
            [action resetForTests];
            [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.4]];
        }
    }
}

- (UIImage *)screenshotWindow:(UIWindow *)window {
    // TODO hardcoded 3.0
    UIGraphicsBeginImageContextWithOptions(window.bounds.size, NO, 3.0f);
    [window drawViewHierarchyInRect:window.bounds afterScreenUpdates:YES];

    UIImage *image = UIGraphicsGetImageFromCurrentImageContext();
    UIGraphicsEndImageContext();

    return image;
}

@end
