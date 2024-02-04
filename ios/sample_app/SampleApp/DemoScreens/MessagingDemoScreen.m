//
//  MessagingDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2024-02-03.
//

#import "MessagingDemoScreen.h"

@import CriticalMoments;
#import "../Utils.h"

@implementation MessagingDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"User Messaging";
        self.infoText = @"These demos show user messaging examples created with the Critical Moments SDK.\n\n • All UI "
                        @"is native and themeable\n • Messages can be triggered by user actions\n • Messages can be "
                        @"conditional based on device state and user behaviour\n • UI, conditions, and trigger are "
                        @"created in config, which can remotely be updated without app store updates";
        self.buttonLink = @"https://docs.criticalmoments.io/actions/actions-overview";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {
    CMDemoAction *appReview = [[CMDemoAction alloc] init];
    appReview.title = @"Ask for App Review";
    appReview.subtitle =
        @"In this example's UI we combine alerts, system app review prompt, and an embedded browser.\n\nSee our docs "
        @"for a guide on using conditional targeting to maximize your app rating including:\n • Check engagment "
        @"history to ensure they are users who are likley to give a positive review\n • Check that it's been a few "
        @"weeks since you last asked\n • Check best other best practices (not low on battery, N days since install, "
        @"not known buggy app version, not old app version, etc)";
    appReview.actionCMActionName = @"app_review_simple";
    [appReview addResetTestTarget:self action:@selector(resetAll)];
    [self addSection:@"App Reviews" withActions:@[ appReview ]];

    CMDemoAction *updateRecommended = [[CMDemoAction alloc] init];
    updateRecommended.title = @"Suggest App Update";
    updateRecommended.subtitle =
        @"Check if app version is very old, or a known buggy app version. If they are, show a modal UI suggesting they "
        @"update in the app store.\n\nInclude a condition such as:\nversionLessThan(app_version, '1.5.2') || "
        @"app_version IN ['buggy_version_1', 'buggy_version_2']";
    updateRecommended.actionCMActionName = @"app_out_of_date";
    [updateRecommended addResetTestTarget:self action:@selector(resetAll)];

    CMDemoAction *outageModal = [[CMDemoAction alloc] init];
    outageModal.title = @"Service Outage or Maintanence Warning";
    outageModal.subtitle =
        @"Server outages happen. CM can let your user know that you're aware, their data is safe, and a fix is coming "
        @"soon. Great messaging in these moments can head off negative reviews, panic about data loss, and many "
        @"support requests.\n\nThese can be remotely added and removed when needed.";
    outageModal.actionCMActionName = @"outage_modal";
    [outageModal addResetTestTarget:self action:@selector(resetAll)];
    [self addSection:@"Modal UI" withActions:@[ updateRecommended, outageModal ]];

    CMDemoAction *productAnnouncement = [[CMDemoAction alloc] init];
    productAnnouncement.title = @"Product Announcements";
    productAnnouncement.subtitle =
        @"Conditionally targeted product announcements can be pushed to the right users:\n\n • New watch app? Announce "
        @"to userw with an Watch\n • Feature announcment only after they upgrade to necessary build\n • Check "
        @"weather and offer a 'Rainy day special'\n • And more: over 100 targeting options";
    productAnnouncement.actionCMActionName = @"product_announcement";
    [productAnnouncement addResetTestTarget:self action:@selector(resetAll)];

    CMDemoAction *tos = [[CMDemoAction alloc] init];
    tos.title = @"Legal Updates";
    tos.subtitle = @"Modal UI combined with embedded browser to inform users of updates, even on old "
                   @"builds.\n\nConditional targeting can check to ensure it is not show if they have already seen it, "
                   @"and not shown if they joined before terms were updated.";
    tos.actionCMActionName = @"legal_update";
    [tos addResetTestTarget:self action:@selector(resetAll)];
    [self addSection:@"Banner UI" withActions:@[ productAnnouncement, tos ]];
}

- (void)resetAll {
    // Modals and alerts
    [Utils.keyWindow.rootViewController dismissViewControllerAnimated:NO completion:nil];

    // Banners
    [CriticalMoments.sharedInstance removeAllBanners];
}

@end
