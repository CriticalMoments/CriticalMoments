//
//  ViewController.h
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-22.
//

#import <UIKit/UIKit.h>

#import "DemoViewContoller.h"

@interface SampleAppCoreViewController : UITabBarController

@property(nonatomic, readonly) CMDemoScreen *demoRoot;

@end
