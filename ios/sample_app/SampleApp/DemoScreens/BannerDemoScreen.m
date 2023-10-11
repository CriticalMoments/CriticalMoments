//
//  BannerDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import "BannerDemoScreen.h"

#import "Utils.h"

@import CriticalMoments;

@interface BannerDemoScreen () <CMBannerActionDelegate>
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

    // Mangement

    CMDemoAction *clearAllBanners = [[CMDemoAction alloc] init];
    clearAllBanners.title = @"Clear all banners";
    clearAllBanners.subtitle = @"Remove all banners from this app";
    clearAllBanners.skipInUiTesting = true;
    clearAllBanners.actionBlock = ^{
      [CMBannerManager.shared removeAllAppWideMessages];
    };

    [self addSection:@"Banner Management"
         withActions:@[
             clearAllBanners,
         ]];

    // Basics

    CMDemoAction *shortBannerAction = [[CMDemoAction alloc] init];
    shortBannerAction.title = @"Short Banner";
    shortBannerAction.subtitle = @"Display a short single line banner, in the default theme";
    shortBannerAction.actionCMActionName = @"short_banner";
    [shortBannerAction addResetTestTarget:self action:@selector(dismissBanners)];

    CMDemoAction *longBannerAction = [[CMDemoAction alloc] init];
    longBannerAction.title = @"Long Banner";
    longBannerAction.subtitle = @"Display a long banner message which will "
                                @"grow to line wrap, in the default theme";
    longBannerAction.actionCMActionName = @"long_banner";
    [longBannerAction addResetTestTarget:self action:@selector(dismissBanners)];

    CMDemoAction *veryLongBannerAction = [[CMDemoAction alloc] init];
    veryLongBannerAction.title = @"Very Long Banner";
    veryLongBannerAction.subtitle = @"Display a very long banner message which "
                                    @"will get truncated, in the default theme";
    veryLongBannerAction.actionCMActionName = @"very_long_banner";
    [veryLongBannerAction addResetTestTarget:self action:@selector(dismissBanners)];

    [self addSection:@"App Wide Banners" withActions:@[ shortBannerAction, longBannerAction, veryLongBannerAction ]];

    // Position

    CMDemoAction *topBanner = [[CMDemoAction alloc] init];
    topBanner.title = @"Top Banner";
    topBanner.subtitle = @"Display a banner on the top of the app, in the default theme";
    topBanner.actionCMActionName = @"top_banner";

    CMDemoAction *swapPosition = [[CMDemoAction alloc] init];
    swapPosition.title = @"Swap banner position";
    swapPosition.subtitle = @"Swap the banner location between the top and bottom.";
    [swapPosition addTarget:self action:@selector(swapBannerPosition)];
    [swapPosition addResetTestTarget:self action:@selector(dismissBanners)];

    [self addSection:@"Banners Position"
         withActions:@[
             topBanner,
             swapPosition,
         ]];

    // Display Options

    CMDemoAction *customThemeBanner = [[CMDemoAction alloc] init];
    customThemeBanner.title = @"Custom Theme Banner";
    customThemeBanner.subtitle = @"Display a banner built from config with custom theme and action";
    customThemeBanner.actionCMActionName = @"custom_theme_banner";
    [customThemeBanner addResetTestTarget:self action:@selector(dismissBanners)];

    CMDemoAction *undismissableBanner = [[CMDemoAction alloc] init];
    undismissableBanner.title = @"Show undismissable banner";
    undismissableBanner.subtitle = @"Show a banner that doesn't have an X to dismiss";
    undismissableBanner.actionCMActionName = @"undismissable_banner";
    [undismissableBanner addResetTestTarget:self action:@selector(dismissBanners)];

    CMDemoAction *singleLineAction = [[CMDemoAction alloc] init];
    singleLineAction.title = @"Show single line banner";
    singleLineAction.subtitle = @"Show a banner that truncates using `maxLineCount` option";
    singleLineAction.actionCMActionName = @"single_line_banner";
    [singleLineAction addResetTestTarget:self action:@selector(dismissBanners)];

    [self addSection:@"Banners Display Options"
         withActions:@[ customThemeBanner, undismissableBanner, singleLineAction ]];

    // Hardcoded

    CMDemoAction *codeBanner = [[CMDemoAction alloc] init];
    codeBanner.title = @"Hardcoded banner";
    codeBanner.subtitle = @"Show a banner using code instead of config. The banner's apearance "
                          @"and action are hardcoded in this sample app.";
    [codeBanner addTarget:self action:@selector(showMessageFromCode)];
    [codeBanner addResetTestTarget:self action:@selector(dismissBanners)];

    [self addSection:@"Banners from Code" withActions:@[ codeBanner ]];
}

- (void)dismissBanners {
    [CMBannerManager.shared removeAllAppWideMessages];
}

- (void)swapBannerPosition {
    if (CMBannerManager.shared.appWideBannerPosition == CMBannerPositionTop) {
        CMBannerManager.shared.appWideBannerPosition = CMBannerPositionBottom;
    } else {
        CMBannerManager.shared.appWideBannerPosition = CMBannerPositionTop;
    }
}

- (void)showMessageFromCode {
    NSString *messageString = @"This banner is created in code instead of config. The same options "
                              @"are available in code if you need them.";
    CMBannerMessage *bannerMessage = [[CMBannerMessage alloc] initWithBody:messageString];
    bannerMessage.actionDelegate = self;
    if (@available(iOS 13, *)) {
        [[CMBannerManager shared] showAppWideMessage:bannerMessage];
    } else {
        [self showAlertWithTitle:@"Not Supported" andBody:@"Banners are not supported on iOS 12 or earlier."];
    }
}

- (void)showAlertWithTitle:(NSString *)title andBody:(NSString *)body {
    UIAlertController *alert = [UIAlertController alertControllerWithTitle:title
                                                                   message:body
                                                            preferredStyle:UIAlertControllerStyleAlert];
    UIAlertAction *defaultAction = [UIAlertAction actionWithTitle:@"OK"
                                                            style:UIAlertActionStyleDefault
                                                          handler:^(UIAlertAction *action){
                                                          }];
    [alert addAction:defaultAction];

    UIViewController *rootVC = Utils.keyWindow.rootViewController;
    [rootVC presentViewController:alert animated:YES completion:nil];
}

#pragma mark CMBannerActionDelegate

- (void)messageAction:(CMBannerMessage *)message {
    NSString *alertMessage = [NSString stringWithFormat:@"Assign an actionDelegate to make this do whatever "
                                                        @"you want!\n\nThe banner you tapped said:\"%@\"",
                                                        message.body];
    [self showAlertWithTitle:@"Banner Tapped" andBody:alertMessage];
}

@end
