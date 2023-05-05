//
//  BannerDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import "BannerDemoScreen.h"

#import "Utils.h"

@import CriticalMoments;

@import Appcore;

@interface BannerDemoScreen () <CMBannerActionDelegate>

@property(nonatomic) NSInteger counter;

@end

@implementation BannerDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Banner Demos";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {

    // Basics

    CMDemoAction *shortBannerAction = [[CMDemoAction alloc] init];
    shortBannerAction.title = @"Show Short Banner";
    shortBannerAction.subtitle =
        @"Display a short single line banner, across entire app";
    [shortBannerAction addTarget:self action:@selector(showShortMessage)];

    CMDemoAction *longBannerAction = [[CMDemoAction alloc] init];
    longBannerAction.title = @"Show Long Banner";
    longBannerAction.subtitle =
        @"Display a long banner message, across entire app";
    [longBannerAction addTarget:self action:@selector(showLongMessage)];

    CMDemoAction *veryLongBannerAction = [[CMDemoAction alloc] init];
    veryLongBannerAction.title = @"Show Very Long Banner";
    veryLongBannerAction.subtitle =
        @"Display a very long banner message, across entire app";
    [veryLongBannerAction addTarget:self action:@selector(showVeryLongMessage)];

    // TODO: remove this -- just for internal testing pre v1
    CMDemoAction *appcoreBannerAction = [[CMDemoAction alloc] init];
    appcoreBannerAction.title = @"Show Banner from Config";
    appcoreBannerAction.subtitle = @"Display a banner built from config";
    [appcoreBannerAction addTarget:self action:@selector(showAppcoreBanner)];

    [self addSection:@"App Wide Banners"
         withActions:@[
             shortBannerAction, longBannerAction, veryLongBannerAction,
             appcoreBannerAction
         ]];

    // Position

    CMDemoAction *swapPosition = [[CMDemoAction alloc] init];
    swapPosition.title = @"Swap banner position";
    swapPosition.subtitle =
        @"Swap the banner location between the top and bottom.";
    [swapPosition addTarget:self action:@selector(swapBannerPosition)];

    [self addSection:@"Banners Position"
         withActions:@[
             swapPosition,
         ]];

    // Mangement

    CMDemoAction *clearAllBanners = [[CMDemoAction alloc] init];
    clearAllBanners.title = @"Clear all banners";
    clearAllBanners.subtitle = @"Remove all banners from this app";
    clearAllBanners.actionBlock = ^{
      [CMBannerManager.shared removeAllAppWideMessages];
    };

    [self addSection:@"Banners Management"
         withActions:@[
             clearAllBanners,
         ]];

    // Display Options

    CMDemoAction *undismissableBanner = [[CMDemoAction alloc] init];
    undismissableBanner.title = @"Show undismissable banner";
    undismissableBanner.subtitle =
        @"Show a banner that doesn't have an X to dismiss";
    [undismissableBanner addTarget:self
                            action:@selector(showUndismissableBanner)];

    CMDemoAction *singleLineAction = [[CMDemoAction alloc] init];
    singleLineAction.title = @"Show single line banner";
    singleLineAction.subtitle =
        @"Show a banner that truncates using `maxLineCount`";
    [singleLineAction addTarget:self action:@selector(showSingleLineMessage)];

    [self addSection:@"Banners Display Options"
         withActions:@[ undismissableBanner, singleLineAction ]];
}

- (void)showShortMessage {
    [self showAppWideBanner:@"Short message"];
}

- (void)showLongMessage {
    [self showAppWideBanner:@"Welcome to critical moments! App wide banners "
                            @"can give your users critical information."];
}

- (void)showVeryLongMessage {
    [self showAppWideBanner:
              @"Welcome to critical moments! App wide banners can give your "
              @"users critical information. This one happens to be really "
              @"really long, and will probably be truncated eventually. It was "
              @"the best of times, it was the worst of times, it was the age "
              @"of wisdom, it was the age of foolishness, it was the epoch of "
              @"belief, it was the epoch of incredulity, it was the season of "
              @"light, it was the season of darkness, it was the spring of "
              @"hope, it was the winter of despair"];
}

- (void)showAppcoreBanner {
    AppcoreInternalDipatchBannerFromGo();
}

- (void)swapBannerPosition {
    if (CMBannerManager.shared.appWideBannerPosition == CMBannerPositionTop) {
        CMBannerManager.shared.appWideBannerPosition = CMBannerPositionBottom;
    } else {
        CMBannerManager.shared.appWideBannerPosition = CMBannerPositionTop;
    }
}

- (void)showUndismissableBanner {
    CMBannerMessage *bannerMessage =
        [[CMBannerMessage alloc] initWithBody:@"You are stuck with me."];
    bannerMessage.actionDelegate = self;
    bannerMessage.showDismissButton = NO;
    [[CMBannerManager shared] showAppWideMessage:bannerMessage];
}

- (void)showSingleLineMessage {
    CMBannerMessage *bannerMessage = [[CMBannerMessage alloc]
        initWithBody:@"This message will truncate after the first line, unlike "
                     @"the default."];
    bannerMessage.actionDelegate = self;
    bannerMessage.maxLineCount = @1;
    [[CMBannerManager shared] showAppWideMessage:bannerMessage];
}

- (void)showAppWideBanner:(NSString *)messageString {
    self.counter += 1;
    NSString *messageStingWithCount = [NSString
        stringWithFormat:@"(%ld) %@", (long)self.counter, messageString];
    CMBannerMessage *bannerMessage =
        [[CMBannerMessage alloc] initWithBody:messageStingWithCount];
    bannerMessage.actionDelegate = self;
    [[CMBannerManager shared] showAppWideMessage:bannerMessage];
}

#pragma mark CMBannerActionDelegate

- (void)messageAction:(CMBannerMessage *)message {
    NSString *alertMessage = [NSString
        stringWithFormat:@"Assign an actionDelegate to make this do whatever "
                         @"you want!\n\nThe banner you tapped said:\"%@\"",
                         message.body];
    UIAlertController *alert = [UIAlertController
        alertControllerWithTitle:@"Banner Tapped"
                         message:alertMessage
                  preferredStyle:UIAlertControllerStyleAlert];
    UIAlertAction *defaultAction =
        [UIAlertAction actionWithTitle:@"OK"
                                 style:UIAlertActionStyleDefault
                               handler:^(UIAlertAction *action){
                               }];
    [alert addAction:defaultAction];

    UIViewController *rootVC = Utils.keyWindow.rootViewController;
    [rootVC presentViewController:alert animated:YES completion:nil];
}

@end
