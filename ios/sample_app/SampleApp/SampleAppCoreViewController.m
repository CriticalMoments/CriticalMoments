//
//  ViewController.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-22.
//

#import "SampleAppCoreViewController.h"

#import "DemoViewContoller.h"
#import "MainDemoScreen.h"
#import "Utils.h"

// TODO_P0
#import "UserNotifications/UserNotifications.h"

#define BANNER_HEIGHT 60.0

@import CriticalMoments;

@interface SampleAppCoreViewController ()

@property(nonatomic, strong) CMDemoScreen *demoRoot;

@end

@implementation SampleAppCoreViewController

- (void)viewDidLoad {
    [super viewDidLoad];

    CMDemoScreen *mainDemoScreen = [[MainDemoScreen alloc] init];
    _demoRoot = mainDemoScreen;
    DemoViewContoller *mainTabRoot = [[DemoViewContoller alloc] initWithDemoScreen:mainDemoScreen];
    UINavigationController *mainTabNav = [[UINavigationController alloc] initWithRootViewController:mainTabRoot];
    mainTabNav.navigationBar.prefersLargeTitles = YES;

    UITabBarItem *mainTabBarItem;
    if (@available(iOS 13.0, *)) {
        UIImageConfiguration *largeImageConfig =
            [UIImageSymbolConfiguration configurationWithScale:UIImageSymbolScaleLarge];
        UIImage *mainTabIconImage = [UIImage systemImageNamed:@"wand.and.rays" withConfiguration:largeImageConfig];
        mainTabBarItem = [[UITabBarItem alloc] initWithTitle:@"Demo" image:mainTabIconImage tag:0];
    } else {
        mainTabBarItem = [[UITabBarItem alloc] initWithTabBarSystemItem:UITabBarSystemItemFeatured tag:0];
    }
    mainTabNav.tabBarItem = mainTabBarItem;
    self.viewControllers = @[ mainTabNav ];

    // Only visible in snapshot test cases where we hack the opacity of the
    // views over this
    if (self.backgroundColor) {
        mainTabNav.view.backgroundColor = self.backgroundColor;
    }

    if (@available(iOS 15.0, *)) {
        UITabBarAppearance *tabAppearance = [[UITabBarAppearance alloc] init];
        [tabAppearance configureWithOpaqueBackground];
        self.tabBar.scrollEdgeAppearance = tabAppearance;
    }

    // TODO_P0: not here!
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    UNAuthorizationOptions opt = UNAuthorizationOptionAlert | UNAuthorizationOptionBadge | UNAuthorizationOptionSound;
    [center requestAuthorizationWithOptions:opt
                          completionHandler:^(BOOL granted, NSError *_Nullable error) {
                            NSLog(@"");
                            // TODO_P0: do we need to trigger scheduling once granted? I think we do!
                            // Maybe
                            // https://developer.apple.com/documentation/usernotifications/unusernotificationcenter/getnotificationsettings(completionhandler:)?language=objc
                          }];
}

- (void)setBackgroundColor:(UIColor *)backgroundColor {
    _backgroundColor = backgroundColor;
    if (self.view) {
        self.viewControllers.firstObject.view.backgroundColor = backgroundColor;
    }
}

@end
