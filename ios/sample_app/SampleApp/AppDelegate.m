//
//  AppDelegate.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-22.
//

#import "AppDelegate.h"

@import CriticalMoments;

#import "Utils.h"

@interface AppDelegate ()

@end

@implementation AppDelegate

- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions {

    return YES;
}

- (id)cmInstance {
    return CriticalMoments.sharedInstance;
}

#pragma mark - UISceneSession lifecycle

- (UISceneConfiguration *)application:(UIApplication *)application
    configurationForConnectingSceneSession:(UISceneSession *)connectingSceneSession
                                   options:(UISceneConnectionOptions *)options API_AVAILABLE(ios(13.0)) {
    // Called when a new scene session is being created.
    // Use this method to select a configuration to create the new scene with.
    return [[UISceneConfiguration alloc] initWithName:@"Default Configuration" sessionRole:connectingSceneSession.role];
}

- (void)application:(UIApplication *)application
    didDiscardSceneSessions:(NSSet<UISceneSession *> *)sceneSessions API_AVAILABLE(ios(13.0)) {
    // Called when the user discards a scene session.
    // If any sessions were discarded while the application was not running,
    // this will be called shortly after
    // application:didFinishLaunchingWithOptions. Use this method to release any
    // resources that were specific to the discarded scenes, as they will not
    // return.
}

- (BOOL)application:(UIApplication *)app
            openURL:(NSURL *)url
            options:(NSDictionary<UIApplicationOpenURLOptionsKey, id> *)options {
    if ([@"critical-moments-sampleapp:main" isEqualToString:url.absoluteString]) {
        // return to the main screen of the app
        UIViewController *rootVC = Utils.keyWindow.rootViewController;
        if ([rootVC isKindOfClass:[UITabBarController class]]) {
            UITabBarController *tab = (UITabBarController *)rootVC;
            rootVC = tab.selectedViewController;
        }
        UINavigationController *navController;
        if ([rootVC isKindOfClass:[UINavigationController class]]) {
            navController = (UINavigationController *)rootVC;
        } else {
            navController = rootVC.navigationController;
        }
        [navController popToRootViewControllerAnimated:YES];
        return YES;
    }
    return NO;
}

@end
