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
    
    //return;;
    
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 0.25 * NSEC_PER_SEC), dispatch_get_main_queue(), ^{
        // TODO better primary -- UIApplication.shared.windows.first(where: { $0.isKeyWindow })?.addSubview(myView)
        // TODO Warning
        NSArray* a = [[UIApplication sharedApplication] windows];
        UIWindow* w = [[[UIApplication sharedApplication] windows] firstObject];
        UIViewController* rootVc = w.rootViewController;
        
        UIView* v = [[CMBannerMessage alloc] initWithBody:@"Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world "];
        //v.backgroundColor = [UIColor greenColor];
        //v.frame = CGRectMake(50,50,50,50);
        v.translatesAutoresizingMaskIntoConstraints = NO;
        v.accessibilityIdentifier = @"banner";
        [w addSubview:v];
        UIView* rootView = rootVc.view;
        
        // Bottom || Top
        if (false) {
            // Bottom
            UIEdgeInsets additionalInset = rootVc.additionalSafeAreaInsets;
            additionalInset.bottom = additionalInset.bottom + BANNER_HEIGHT;
            rootVc.additionalSafeAreaInsets = additionalInset;
            
            NSArray<NSLayoutConstraint*>* constraints = @[
                
                // position below the window and to the edges
                [v.topAnchor constraintEqualToAnchor:rootView.layoutMarginsGuide.bottomAnchor],
                [v.leftAnchor constraintEqualToAnchor:w.leftAnchor],
                [v.rightAnchor constraintEqualToAnchor:w.rightAnchor],
                [v.bottomAnchor constraintEqualToAnchor:w.bottomAnchor],
            ];
            
            [NSLayoutConstraint activateConstraints:constraints];
        } else {
            // Top
            UIEdgeInsets additionalInset = rootVc.additionalSafeAreaInsets;
            additionalInset.top = additionalInset.top + BANNER_HEIGHT;
            rootVc.additionalSafeAreaInsets = additionalInset;
            
            NSArray<NSLayoutConstraint*>* constraints = @[
                
                // position below the window and to the edges
                [v.topAnchor constraintEqualToAnchor:w.topAnchor],
                [v.leftAnchor constraintEqualToAnchor:w.leftAnchor],
                [v.rightAnchor constraintEqualToAnchor:w.rightAnchor],
                [v.bottomAnchor constraintEqualToAnchor:rootView.layoutMarginsGuide.topAnchor],
            ];
            
            [NSLayoutConstraint activateConstraints:constraints];
        }
        
        [rootView setNeedsLayout];
    });
}


@end
