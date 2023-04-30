//
//  ViewController.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-22.
//

#import "ViewController.h"

#import "DemoViewContoller.h"
#import "MainDemoScreen.h"

#define BANNER_HEIGHT 60.0

@import CriticalMoments;

@interface ViewController ()

@end

@implementation ViewController

- (void)viewDidLoad {
    [super viewDidLoad];

    CMDemoScreen *mainDemoScreen = [[MainDemoScreen alloc] init];
    DemoViewContoller *mainTabRoot =
        [[DemoViewContoller alloc] initWithDemoScreen:mainDemoScreen];
    UINavigationController *mainTabNav =
        [[UINavigationController alloc] initWithRootViewController:mainTabRoot];
    mainTabNav.navigationBar.prefersLargeTitles = YES;

    UIImageConfiguration *largeImageConfig = [UIImageSymbolConfiguration
        configurationWithScale:UIImageSymbolScaleLarge];
    UIImage *mainTabIconImage = [UIImage systemImageNamed:@"wand.and.rays"
                                        withConfiguration:largeImageConfig];
    UITabBarItem *mainTabBarItem =
        [[UITabBarItem alloc] initWithTitle:@"Demo"
                                      image:mainTabIconImage
                                        tag:0];
    mainTabNav.tabBarItem = mainTabBarItem;
    self.viewControllers = @[ mainTabNav ];

    if (@available(iOS 15.0, *)) {
        UITabBarAppearance *tabAppearance = [[UITabBarAppearance alloc] init];
        [tabAppearance configureWithOpaqueBackground];
        self.tabBar.scrollEdgeAppearance = tabAppearance;
    }
}

@end
