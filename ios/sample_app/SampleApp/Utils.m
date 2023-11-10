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

static NSArray<NSURL *> *writeableDirs = nil;
static NSURL *bundleUrl = nil;

+ (void)createTestFileUrls {
    NSURL *appSupportDir = [[NSFileManager.defaultManager URLsForDirectory:NSApplicationSupportDirectory
                                                                 inDomains:NSUserDomainMask] lastObject];
    NSURL *cachesDir = [[NSFileManager.defaultManager URLsForDirectory:NSCachesDirectory
                                                             inDomains:NSUserDomainMask] lastObject];
    NSURL *docsDir = [[NSFileManager.defaultManager URLsForDirectory:NSDocumentDirectory
                                                           inDomains:NSUserDomainMask] lastObject];
    NSURL *tempDirUrl = [NSURL fileURLWithPath:NSTemporaryDirectory()];

    writeableDirs = @[ tempDirUrl, cachesDir, appSupportDir, docsDir ];
    bundleUrl = [[NSBundle mainBundle] bundleURL];
}

+ (BOOL)verifyTestFileUrls {
    for (NSURL *dir in writeableDirs) {
        if ([dir.absoluteString hasPrefix:bundleUrl.absoluteString]) {
            return false;
        }
    }

    return true;
}

@end
