//
//  MainTabViewContoller.h
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-23.
//

#import <UIKit/UIKit.h>
#import "CMDemoScreen.h"

NS_ASSUME_NONNULL_BEGIN

@interface DemoViewContoller : UITableViewController

-(instancetype)init NS_UNAVAILABLE;
-(instancetype)initWithDemoScreen:(CMDemoScreen*)screen;

/*
 demo manager class:
  - add a ton of little action methods that "do the demo"
  - Big list of menu items: title, subtitle, P2 image, target method (like UIButton's), data (custom data target knows how to use)
  - P2: sections
  - a way to make an app structure: menus>sections. And actions
 */

@end

NS_ASSUME_NONNULL_END
