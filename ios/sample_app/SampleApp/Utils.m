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

/// Parse hex codes in format #ffffff to UIColor
+ (UIColor *)colorFromHexString:(NSString *)hexString {
    if (hexString.length != 7) {
        return nil;
    }

    unsigned int parsed = 0;
    NSScanner *scanner = [NSScanner scannerWithString:hexString];

    if ([hexString hasPrefix:@"#"]) {
        [scanner setScanLocation:1];
    } else {
        return nil;
    }
    bool scannedHex = [scanner scanHexInt:&parsed];
    if (!scannedHex || ![scanner isAtEnd]) {
        return nil;
    }

    CGFloat red = ((parsed & 0xff0000) >> 16) / 255.0;
    CGFloat green = ((parsed & 0x00ff00) >> 8) / 255.0;
    CGFloat blue = (parsed & 0x0000ff) / 255.0;

    return [[UIColor alloc] initWithRed:red green:green blue:blue alpha:1.0];
}

+ (NSError *)deleteDatabase {
    // This is only for the demo app. You really really shouldn't emulate this in a client app. This code is not
    // guaranteed to work over time, not is deleting the database file a good idea.
    NSURL *appSupportDir = [[NSFileManager.defaultManager URLsForDirectory:NSApplicationSupportDirectory
                                                                 inDomains:NSUserDomainMask] lastObject];

    // Set the data directory to applicationSupport/critical_moments_data
    NSError *error;
    NSURL *criticalMomentsDataDir = [appSupportDir URLByAppendingPathComponent:@"critical_moments_data"];

    BOOL success = [NSFileManager.defaultManager removeItemAtURL:criticalMomentsDataDir error:&error];
    if (!success || error) {
        // Errors okay. Fresh start won't have a cache yet.
        return error;
    }
    return nil;
}

@end
