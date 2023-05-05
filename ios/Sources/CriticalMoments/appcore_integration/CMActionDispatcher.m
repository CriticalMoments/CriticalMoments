//
//  CMActionDispatcher.m
//
//
//  Created by Steve Cosman on 2023-05-05.
//

#import "CMActionDispatcher.h"

#import "../messaging/CMBannerManager.h"

@import Appcore;

@interface CMActionDispatcher () <AppcoreLibActionDispatcher>
@end

@implementation CMActionDispatcher

static CMActionDispatcher *sharedInstance = nil;

+ (CMActionDispatcher *)shared {
    // avoid lock if we can
    if (sharedInstance) {
        return sharedInstance;
    }

    @synchronized(CMActionDispatcher.class) {
        if (!sharedInstance) {
            sharedInstance = [[self alloc] init];
        }

        return sharedInstance;
    }
}

+ (void)registerWithAppcore {
    AppcoreRegisterActionDispatcher([CMActionDispatcher shared]);
}

- (void)showBanner:(DatamodelBannerAction *_Nullable)bannerAction {
    if (!bannerAction) {
        return;
    }

    CMBannerMessage *bannerMessage =
        [[CMBannerMessage alloc] initWithAppcoreDataModel:bannerAction];
    [[CMBannerManager shared] showAppWideMessage:bannerMessage];
}

@end
