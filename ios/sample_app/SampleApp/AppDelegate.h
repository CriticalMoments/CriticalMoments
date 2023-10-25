//
//  AppDelegate.h
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-22.
//

#import <UIKit/UIKit.h>

@import CriticalMoments;

@interface AppDelegate : UIResponder <UIApplicationDelegate>

@property(strong, nonatomic) UIWindow *window;

- (CriticalMoments *)cmInstance;

@end
