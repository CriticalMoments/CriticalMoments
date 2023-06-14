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

    CMBannerMessage *bannerMessage = [[CMBannerMessage alloc] initWithAppcoreDataModel:banner];

    // TODO: main thread?
    [[CMBannerManager shared] showAppWideMessage:bannerMessage];
    return YES;
}

- (BOOL)showAlert:(DatamodelAlertAction *_Nullable)alertDataModel
            error:(NSError *_Nullable __autoreleasing *_Nullable)error {
    if (!alertDataModel) {
        *error = [NSError errorWithDomain:@"CMIOS" code:4565684 userInfo:nil];
        return NO;
    }
    CMAlert *alert = [[CMAlert alloc] initWithAppcoreDataModel:alertDataModel];
    [alert showAlert];
    return YES;
}

- (BOOL)showLink:(DatamodelLinkAction *)link error:(NSError *_Nullable __autoreleasing *)error {
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
}

- (BOOL)openLinkInEmbeddedBrowser:(NSURL *)url {
    SFSafariViewController *safariVc = [[SFSafariViewController alloc] initWithURL:url];
    UIViewController *rootVc = CMUtils.keyWindow.rootViewController;
    if (!safariVc || !rootVc) {
        return NO;
    }
    [rootVc presentViewController:safariVc animated:YES completion:nil];
    return YES;
}

@end
