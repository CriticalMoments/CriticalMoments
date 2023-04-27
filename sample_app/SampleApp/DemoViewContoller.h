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

@end

NS_ASSUME_NONNULL_END
