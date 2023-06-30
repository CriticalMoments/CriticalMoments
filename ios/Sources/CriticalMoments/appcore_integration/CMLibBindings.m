//
//  CMActionDispatcher.m
//
//
//  Created by Steve Cosman on 2023-05-05.
//

#import "CMLibBindings.h"

#import "../messaging/CMAlert.h"
#import "../messaging/CMAlert_private.h"
#import "../messaging/CMBannerManager.h"
#import "../messaging/CMBannerMessage.h"
#import "../messaging/CMBannerMessage_private.h"
#import "../messaging/CMModalViewController.h"
#import "../themes/CMTheme.h"
#import "../themes/CMTheme_private.h"
#import "../utils/CMUtils.h"

@import Appcore;

@import SafariServices;
@import StoreKit;

@interface CMLibBindings () <AppcoreLibBindings>
@end

@implementation CMLibBindings

static CMLibBindings *sharedInstance = nil;

+ (CMLibBindings *)shared {
    // avoid lock if we can
    if (sharedInstance) {
        return sharedInstance;
    }

    @synchronized(CMLibBindings.class) {
        if (!sharedInstance) {
            sharedInstance = [[self alloc] init];
        }

        return sharedInstance;
    }
}

+ (void)registerWithAppcore {
    [AppcoreSharedAppcore() registerLibraryBindings:[CMLibBindings shared]];
}

#pragma mark AppcoreLibBindings

- (BOOL)setDefaultTheme:(DatamodelTheme *_Nullable)actheme error:(NSError *_Nullable __autoreleasing *_Nullable)error {
    if (!actheme) {
        *error = [NSError errorWithDomain:@"CMIOS" code:73923755 userInfo:nil];
        return NO;
    }

    CMTheme *theme = [CMTheme themeFromAppcoreTheme:actheme];
    if (!theme) {
        *error = [NSError errorWithDomain:@"CMIOS" code:81263223 userInfo:nil];
        return NO;
    }
    [CMTheme setCurrentTheme:theme];

    return YES;
}

// TODO test case
- (BOOL)showBanner:(DatamodelBannerAction *_Nullable)banner error:(NSError *_Nullable __autoreleasing *_Nullable)error {
    if (!banner) {
        *error = [NSError errorWithDomain:@"CMIOS" code:92739238 userInfo:nil];
        return NO;
    }

    if (@available(iOS 13, *)) {
        // TODO: main thread?
        CMBannerMessage *bannerMessage = [[CMBannerMessage alloc] initWithAppcoreDataModel:banner];
        [[CMBannerManager shared] showAppWideMessage:bannerMessage];
    } else {
        NSLog(@"CriticalMoments: Banner messages are only supported on iOS 13 or newer.");
        *error = [NSError errorWithDomain:@"CMIOS: Banner not supported on iOS version" code:87155467 userInfo:nil];
        return NO;
    }

    // TODO: what is the bool return here?
    return YES;
}

- (BOOL)showAlert:(DatamodelAlertAction *_Nullable)alertDataModel
            error:(NSError *_Nullable __autoreleasing *_Nullable)error {
    if (!alertDataModel) {
        *error = [NSError errorWithDomain:@"CMIOS" code:4565684 userInfo:nil];
        return NO;
    }

    // TODO no dispatch
    CMAlert *alert = [[CMAlert alloc] initWithAppcoreDataModel:alertDataModel];
    [alert showAlert];

    // TODO: what is the bool return here?
    return YES;
}

- (BOOL)showLink:(DatamodelLinkAction *)link error:(NSError *_Nullable __autoreleasing *)error {
    // TODO no dispatch to main

    NSURL *url = [NSURL URLWithString:link.urlString];
    if (!url || !url.scheme) {
        *error = [NSError errorWithDomain:@"CMIOS" code:72937634 userInfo:nil];
        return NO;
    }

    BOOL isWebLink = [@"http" isEqualToString:url.scheme] || [@"https" isEqualToString:url.scheme];
    if (link.useEmbeddedBrowser && isWebLink) {
        BOOL success = [self openLinkInEmbeddedBrowser:url];
        if (success) {
            return YES;
        }
    }

    [UIApplication.sharedApplication openURL:url options:@{} completionHandler:nil];
    // TODO: what is the bool return here?
    return YES;
}

- (BOOL)showReviewPrompt:(NSError *_Nullable __autoreleasing *)error {
    dispatch_async(dispatch_get_main_queue(), ^{
      if (@available(iOS 14.0, *)) {
          UIWindowScene *scene = [CMUtils keyWindow].windowScene;
          if (scene) {
              [SKStoreReviewController requestReviewInScene:scene];
          }
      } else {
          [SKStoreReviewController requestReview];
      }
    });

    // TODO no returns
}

- (BOOL)showModal:(DatamodelModalAction *_Nullable)modal error:(NSError *_Nullable __autoreleasing *_Nullable)error {
    dispatch_async(dispatch_get_main_queue(), ^{
      CMModalViewController *sheetVc = [[CMModalViewController alloc] initWithDatamodel:modal];
      [CMUtils.topViewController presentViewController:sheetVc animated:YES completion:nil];
    });

    return NO;
}

- (BOOL)openLinkInEmbeddedBrowser:(NSURL *)url {
    SFSafariViewController *safariVc = [[SFSafariViewController alloc] initWithURL:url];
    UIViewController *topController = CMUtils.topViewController;
    if (!safariVc || !topController) {
        return NO;
    }
    [topController presentViewController:safariVc animated:YES completion:nil];

    return YES;
}

@end
