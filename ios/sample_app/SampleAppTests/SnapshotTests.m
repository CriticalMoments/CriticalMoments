//
//  SnapshotTests.m
//  SampleAppTests
//
//  Created by Steve Cosman on 2023-05-11.
//

#import <XCTest/XCTest.h>

#import "SampleAppTests-Swift.h"

#import "SampleAppCoreViewController.h"
#import "Utils.h"

#import <sys/utsname.h>

@interface SnapshotTests : XCTestCase

@end

@implementation SnapshotTests

- (void)setUp {
    [super setUp];

    self.continueAfterFailure = true;
}

- (void)tearDown {
    [super tearDown];
}

- (void)testScreenshotAllSampleAppFeatures {
    [self testAllFeaturesDarkMode:false];
    [self testAllFeaturesDarkMode:true];
}

- (void)testAllFeaturesDarkMode:(bool)darkMode {
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

    // set dark mode
    window.overrideUserInterfaceStyle = darkMode ? UIUserInterfaceStyleDark : UIUserInterfaceStyleLight;

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
                // Take screenshot and diff
                UIImage *screenshot = [self screenshotWindow:[Utils keyWindow]];
                NSString *testName = [self buildNameForDeviceAndAction:action.title withWindow:[Utils keyWindow]];
                CMSnapshotWrapper *w = [[CMSnapshotWrapper alloc] init];
                [w assertSnapshotImageOf:screenshot named:testName];
            }

            // reset state for next test, and give it time to render
            [action resetForTests];
            [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.2]];
        }
    }
}

- (UIImage *)screenshotWindow:(UIWindow *)window {
    UIGraphicsBeginImageContextWithOptions(window.bounds.size, NO, window.screen.scale);
    [window drawViewHierarchyInRect:window.bounds afterScreenUpdates:YES];
    UIImage *image = UIGraphicsGetImageFromCurrentImageContext();
    UIGraphicsEndImageContext();
    return image;
}

- (NSString *)buildNameForDeviceAndAction:(NSString *)actionName withWindow:(UIWindow *)window {
    struct utsname systemInfo;
    uname(&systemInfo);
    NSString *deviceModel = [NSString stringWithCString:systemInfo.machine encoding:NSUTF8StringEncoding];

    if ([@[ @"arm64", @"i386", @"x86_64" ] containsObject:deviceModel]) {
        // This is a simulator. They don't return a model_version_number
        deviceModel = [NSString stringWithFormat:@"%s-Simulator", getenv("SIMULATOR_MODEL_IDENTIFIER")];
    }

    NSString *systemOsName = UIDevice.currentDevice.systemName;
    if ([@"iOS" isEqualToString:systemOsName] && UIDevice.currentDevice.userInterfaceIdiom == UIUserInterfaceIdiomPad) {
        systemOsName = @"iPadOS";
    }

    NSString *darkMode = @"";
    if (window.rootViewController.traitCollection.userInterfaceStyle == UIUserInterfaceStyleDark) {
        darkMode = @"[darkmode]";
    }

    return [NSString stringWithFormat:@"%@%@-%dx%d@%.02f-%@--%@", deviceModel, darkMode, (int)window.bounds.size.width,
                                      (int)window.bounds.size.height, window.screen.scale, systemOsName, actionName];
}

@end
