//
//  Utils.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import "Utils.h"

@implementation Utils

+ (UIWindow *)keyWindow {
    UIWindow *keyWindow = [[[UIApplication sharedApplication] windows] firstObject];
    for (UIWindow *w in [[UIApplication sharedApplication] windows]) {
        if (w.isKeyWindow) {
            keyWindow = w;
            break;
        }
    }
    return keyWindow;
}

+ (UINavigationController *)appNavControl {
    UINavigationController *navController;
    UIViewController *rootVC = Utils.keyWindow.rootViewController;
    if ([rootVC isKindOfClass:[UITabBarController class]]) {
        UITabBarController *tab = (UITabBarController *)rootVC;
        rootVC = tab.selectedViewController;
    }
    if ([rootVC isKindOfClass:[UINavigationController class]]) {
        navController = (UINavigationController *)rootVC;
    } else {
        navController = rootVC.navigationController;
    }
    return navController;
}

@end
