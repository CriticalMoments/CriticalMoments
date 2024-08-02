//
//  CMActionDispatcher.m
//
//
//  Created by Steve Cosman on 2023-05-05.
//

#import "CMLibBindings.h"

#import "../CriticalMoments_private.h"
#import "../background/CMBackgroundHandler.h"
#import "../messaging/CMAlert.h"
#import "../messaging/CMAlert_private.h"
#import "../messaging/CMBannerManager.h"
#import "../messaging/CMBannerMessage.h"
#import "../messaging/CMBannerMessage_private.h"
#import "../messaging/CMModalViewController.h"
#import "../notifications/CMNotificationHandler.h"
#import "../themes/CMTheme.h"
#import "../themes/CMTheme_private.h"
#import "../utils/CMUtils.h"

@import Appcore;

@import SafariServices;
@import StoreKit;

@interface CMLibBindings ()
@property(nonatomic, weak) CriticalMoments *cm; // weak to avoid circular reference
@end

@implementation CMLibBindings

#pragma mark AppcoreLibBindings

- (instancetype)initWithCM:(CriticalMoments *)cm {
    self = [super init];
    if (self) {
        _cm = cm;
    }
    return self;
}

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

    [self.cm setTheme:theme];
    return YES;
}

- (BOOL)setDefaultThemeByLibaryThemeName:(NSString *_Nullable)themeName
                                   error:(NSError *_Nullable __autoreleasing *_Nullable)error {
    CMTheme *theme = [CMTheme libaryThemeByName:themeName];
    if (!theme) {
        *error = [NSError errorWithDomain:@"CMIOS" code:902834902384 userInfo:nil];
        return NO;
    }

    [self.cm setTheme:theme];
    return YES;
}

- (BOOL)showBanner:(DatamodelBannerAction *)banner
        actionName:(NSString *)actionName
             error:(NSError *_Nullable __autoreleasing *)error {
    if (!banner) {
        *error = [NSError errorWithDomain:@"CMIOS" code:92739238 userInfo:nil];
        return NO;
    }

    if (@available(iOS 13, *)) {
        dispatch_async(dispatch_get_main_queue(), ^{
          CMBannerMessage *bannerMessage = [[CMBannerMessage alloc] initWithAppcoreDataModel:banner];
          bannerMessage.completionEventSender = self.cm;
          if (actionName.length > 0) {
              bannerMessage.bannerName = actionName;
          }
          [[CMBannerManager shared] showAppWideMessage:bannerMessage];
        });
    } else {
        NSLog(@"CriticalMoments: Banner messages are only supported on iOS 13 or newer.");
        *error = [NSError errorWithDomain:@"CMIOS: Banner not supported on iOS version" code:87155467 userInfo:nil];
        return NO;
    }

    return YES;
}

- (BOOL)showAlert:(DatamodelAlertAction *)alertDataModel
       actionName:(NSString *)actionName
            error:(NSError *_Nullable __autoreleasing *)error {
    if (!alertDataModel) {
        *error = [NSError errorWithDomain:@"CMIOS" code:4565684 userInfo:nil];
        return NO;
    }

    dispatch_async(dispatch_get_main_queue(), ^{
      CMAlert *alert = [[CMAlert alloc] initWithAppcoreDataModel:alertDataModel];
      if (actionName.length > 0) {
          alert.alertName = actionName;
      }
      alert.completionEventSender = self.cm;
      [alert showAlert];
    });

    return YES;
}

- (BOOL)showLink:(DatamodelLinkAction *)link error:(NSError *_Nullable __autoreleasing *)error {
    NSURL *url = [NSURL URLWithString:link.urlString];
    if (!url || !url.scheme) {
        *error = [NSError errorWithDomain:@"CMIOS" code:72937634 userInfo:nil];
        return NO;
    }

    dispatch_async(dispatch_get_main_queue(), ^{
      BOOL isWebLink = [@"http" isEqualToString:url.scheme] || [@"https" isEqualToString:url.scheme];
      if (link.useEmbeddedBrowser && isWebLink) {
          [self openLinkInEmbeddedBrowser:url];
      } else {
          [UIApplication.sharedApplication openURL:url options:@{} completionHandler:nil];
      }
    });
    return YES;
}

- (BOOL)showReviewPrompt:(NSError *_Nullable __autoreleasing *)error {
    __block CriticalMoments *blockCM = self.cm;
    dispatch_async(dispatch_get_main_queue(), ^{
      if (@available(iOS 14.0, *)) {
          UIWindowScene *scene = [CMUtils keyWindow].windowScene;
          if (scene) {
              [SKStoreReviewController requestReviewInScene:scene];
          }
      } else {
          [SKStoreReviewController requestReview];
      }

      [blockCM sendEvent:@"system_app_review_requested"];
    });

    return YES;
}

- (BOOL)showModal:(DatamodelModalAction *)modal
       actionName:(NSString *)actionName
            error:(NSError *_Nullable __autoreleasing *)error {
    dispatch_async(dispatch_get_main_queue(), ^{
      CMModalViewController *sheetVc = [[CMModalViewController alloc] initWithDatamodel:modal];
      if (actionName.length > 0) {
          sheetVc.modalName = actionName;
      }
      sheetVc.completionEventSender = self.cm;
      [CMUtils.topViewController presentViewController:sheetVc animated:YES completion:nil];
    });

    return YES;
}

- (BOOL)updateNotificationPlan:(AppcoreNotificationPlan *_Nullable)notifPlan
                         error:(NSError *_Nullable __autoreleasing *_Nullable)error {
    if ([self.cm userNotificationsDisabled]) {
        return YES;
    }

    dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
      [CMNotificationHandler updateNotificationPlan:notifPlan];

      [self.cm.backgroundHandler scheduleBackgroundTaskAtEpochTime:notifPlan.earliestBgCheckTimeEpochSeconds];
    });
    return YES;
}

- (BOOL)canOpenURL:(NSString *_Nullable)urlString {
    NSURL *url = [NSURL URLWithString:urlString];
    if (url) {
        return [UIApplication.sharedApplication canOpenURL:url];
    }
    return NO;
}

- (NSString *_Nonnull)appVersion {
    NSString *appVersion = [NSBundle.mainBundle objectForInfoDictionaryKey:@"CFBundleShortVersionString"];
    return appVersion;
}

- (NSString *_Nonnull)cmVersion {
    return CM_LIB_VERSION_NUMBER_STRING;
}

- (void)openLinkInEmbeddedBrowser:(NSURL *)url {
    dispatch_async(dispatch_get_main_queue(), ^{
      SFSafariViewController *safariVc = [[SFSafariViewController alloc] initWithURL:url];
      UIViewController *topController = CMUtils.topViewController;
      if (!safariVc || !topController) {
          return;
      }
      [topController presentViewController:safariVc animated:YES completion:nil];
    });
}

// Only used in testing, so while this could break over time, CI will catch it and it
// won't impact production apps.
- (BOOL)isTestBuild {
    bool testEnv = [[[NSProcessInfo processInfo] environment] objectForKey:@"XCTestConfigurationFilePath"] != nil;
    if (testEnv) {
        return true;
    }

    return false;
}

@end
