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
    AppcoreAppcore* appcore = AppcoreSharedAppcore();
    [AppcoreSharedAppcore() registerLibraryBindings:[CMLibBindings shared]];
}

#pragma mark AppcoreLibBindings

- (BOOL)setDefaultTheme:(DatamodelTheme * _Nullable)theme error:(NSError * _Nullable __autoreleasing * _Nullable)error {
    // TODO
    return YES;
}

- (BOOL)showBanner:(DatamodelBannerAction * _Nullable)banner error:(NSError * _Nullable __autoreleasing * _Nullable)error {
    if (!banner) {
        *error = [NSError errorWithDomain:@"CMIOS" code:92739238 userInfo:nil];
        return;
    }

    CMBannerMessage *bannerMessage =
        [[CMBannerMessage alloc] initWithAppcoreDataModel:banner];
    
    if (banner.customThemeName.length > 0) {
        CMTheme* customTheme = [CMTheme namedThemeFromConfig:banner.customThemeName];
        bannerMessage.customTheme = customTheme;
    }
    
    // TODO: main thread?
    [[CMBannerManager shared] showAppWideMessage:bannerMessage];
}


@end
