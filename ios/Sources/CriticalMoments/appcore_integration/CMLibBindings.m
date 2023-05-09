//
//  CMActionDispatcher.m
//
//
//  Created by Steve Cosman on 2023-05-05.
//

#import "CMLibBindings.h"

#import "../messaging/CMBannerManager.h"
#import "../messaging/CMBannerMessage.h"
#import "../messaging/CMBannerMessage_private.h"
#import "../themes/CMTheme.h"
#import "../themes/CMTheme_private.h"

@import Appcore;

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
    AppcoreAppcore *appcore = AppcoreSharedAppcore();
    [AppcoreSharedAppcore() registerLibraryBindings:[CMLibBindings shared]];
}

#pragma mark AppcoreLibBindings

// TODO test case
- (BOOL)setDefaultTheme:(DatamodelTheme *_Nullable)actheme
                  error:(NSError *_Nullable __autoreleasing *_Nullable)error {
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
- (BOOL)showBanner:(DatamodelBannerAction *_Nullable)banner
             error:(NSError *_Nullable __autoreleasing *_Nullable)error {
    if (!banner) {
        *error = [NSError errorWithDomain:@"CMIOS" code:92739238 userInfo:nil];
        return;
    }

    CMBannerMessage *bannerMessage =
        [[CMBannerMessage alloc] initWithAppcoreDataModel:banner];

    // TODO: main thread?
    [[CMBannerManager shared] showAppWideMessage:bannerMessage];
}

@end
