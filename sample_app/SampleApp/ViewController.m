//
//  ViewController.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-22.
//

#import "ViewController.h"

#import "MainTabViewContoller.h"

#define BANNER_HEIGHT 60.0

@import CriticalMomentsObjc;

@interface ViewController ()

@end

@implementation ViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    
    MainTabViewContoller* mainTab = [[MainTabViewContoller alloc] init];
    UITabBarItem* mainTabBarItem = [[UITabBarItem alloc] initWithTabBarSystemItem:UITabBarSystemItemTopRated tag:0];
    mainTab.tabBarItem = mainTabBarItem;
    //mainTab.title = @"Examples";
    //mainTab.tabBarItem.image = [[UIImage alloc] init];
    self.viewControllers = @[mainTab];
    
    UITabBarAppearance* tabAppearance = [[UITabBarAppearance alloc] init];
    [tabAppearance configureWithOpaqueBackground];
    self.tabBar.scrollEdgeAppearance = tabAppearance;
    
    
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 0.25 * NSEC_PER_SEC), dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        CMBannerMessage* message = [[CMBannerMessage alloc] initWithBody:@"Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world "];
        [[CMBannerManager sharedInstance] showAppWideMessage:message];
    });
    
    
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 1.25 * NSEC_PER_SEC), dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        CMBannerMessage* message = [[CMBannerMessage alloc] initWithBody:@"message two, ya ba do"];
        [[CMBannerManager sharedInstance] showAppWideMessage:message];
    });

}


@end
