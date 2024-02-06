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
        self.infoText =
            @"For these demos:\n\n • All UI "
            @"is native and themeable\n • Messages can be triggered by user actions\n • Messages can be "
            @"conditional based on device state and user behaviour history\n • UI, conditions, and triggers are "
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
        @"for a guide on using conditional targeting to maximize your app rating by asking users to rate when they are "
        @"likley to give a positive review (not low battery, not recently installed, explored key features, etc).";
    appReview.actionCMActionName = @"app_review_simple";
    [appReview addResetTestTarget:self action:@selector(resetAll)];
    [self addSection:@"App Reviews" withActions:@[ appReview ]];

    CMDemoAction *updateRecommended = [[CMDemoAction alloc] init];
    updateRecommended.title = @"Suggest App Update";
    updateRecommended.subtitle = @"Check if the app version is very old or a known buggy app version; if so, suggest "
                                 @"they update.\n\nExample condition:\nversionLessThan(app_version, '1.5.2') ||\n"
                                 @"app_version IN ['buggy_version_1', 'buggy_version_2']";
    updateRecommended.actionCMActionName = @"app_out_of_date";
    [updateRecommended addResetTestTarget:self action:@selector(resetAll)];

    CMDemoAction *outageModal = [[CMDemoAction alloc] init];
    outageModal.title = @"Service Outage or Maintanence Warning";
    outageModal.subtitle = @"During an outage, let your user know that you're aware and working on a fix."
                           @" Messaging in these moments can reduce negative reviews, panic about data loss, and "
                           @"support request volume.\n\nThese can be remotely added during outages, and removed after.";
    outageModal.actionCMActionName = @"outage_modal";
    [outageModal addResetTestTarget:self action:@selector(resetAll)];
    [self addSection:@"Modal UI" withActions:@[ updateRecommended, outageModal ]];

    CMDemoAction *productAnnouncement = [[CMDemoAction alloc] init];
    productAnnouncement.title = @"Announcements";
    productAnnouncement.snapshotTitle = @"Product Announcements";
    productAnnouncement.subtitle = @"Examples of conditionally targeted announcements:\n • New watch app? Announce "
                                   @"to users with an Watch\n • Check local "
                                   @"weather and offer a 'Rainy day special'\n • Black friday discount banner appears "
                                   @"at midnight\n • And more: over 100 targeting options";
    productAnnouncement.actionCMActionName = @"product_announcement";
    [productAnnouncement addResetTestTarget:self action:@selector(resetAll)];

    CMDemoAction *tos = [[CMDemoAction alloc] init];
    tos.title = @"Handle the Unexpected";
    tos.snapshotTitle = @"Legal Updates";
    tos.subtitle = @"Remotely push updates to the appriopiate users, even for unexpected/unplanned changes.\n\nExample "
                   @"condition: locale_country_code IN ['DE', 'FR', 'RO'...]";
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
