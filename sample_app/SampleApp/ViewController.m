//
//  ViewController.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-22.
//

#import "ViewController.h"

#import "DemoViewContoller.h"
#import "BannerDemoScreen.h"

#define BANNER_HEIGHT 60.0

@import CriticalMomentsObjc;

@interface ViewController ()

@end

@implementation ViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    
    CMDemoScreen* bannerScreen = [[BannerDemoScreen alloc] init];
    DemoViewContoller* mainTabRoot = [[DemoViewContoller alloc] initWithDemoScreen:bannerScreen];
    UINavigationController* mainTabNav = [[UINavigationController alloc] initWithRootViewController:mainTabRoot];
    mainTabNav.navigationBar.prefersLargeTitles = YES;
    
    UIImageConfiguration *largeImageConfig = [UIImageSymbolConfiguration configurationWithScale:UIImageSymbolScaleLarge];
    UIImage* mainTabIconImage = [UIImage systemImageNamed:@"wand.and.rays" withConfiguration:largeImageConfig];
    UITabBarItem* mainTabBarItem = [[UITabBarItem alloc] initWithTitle:@"Demo" image:mainTabIconImage tag:0];
    mainTabNav.tabBarItem = mainTabBarItem;
    self.viewControllers = @[mainTabNav];
    
    UITabBarAppearance* tabAppearance = [[UITabBarAppearance alloc] init];
    [tabAppearance configureWithOpaqueBackground];
    self.tabBar.scrollEdgeAppearance = tabAppearance;
    
    /*dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 0.25 * NSEC_PER_SEC), dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        CMBannerMessage* message = [[CMBannerMessage alloc] initWithBody:@"Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world "];
        message = [[CMBannerMessage alloc] initWithBody:@"short msg"];
        message.actionDelegate = self;
        [[CMBannerManager sharedInstance] showAppWideMessage:message];
    });
    
    
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 1.25 * NSEC_PER_SEC), dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        CMBannerMessage* message = [[CMBannerMessage alloc] initWithBody:@"message two, ya ba do. message two, ya ba do message two, sdf sdf  sdf sdfsf sdf sdya ba do message two, ya ba do. message two, ya ba do message two, ya ba do"];
        //message.maxLines = @4;
        //message.showDismissButton = NO;
        message.actionDelegate = self;
        [[CMBannerManager sharedInstance] showAppWideMessage:message];
        [CMBannerManager sharedInstance].appWideBannerPosition = CMAppWideBannerPositionTop;
    });*/
}

@end
