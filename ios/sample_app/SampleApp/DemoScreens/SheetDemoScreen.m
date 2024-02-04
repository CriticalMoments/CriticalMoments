//
//  SheetDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-06-14.
//

#import "SheetDemoScreen.h"

#import "Utils.h"

@import CriticalMoments;

@implementation SheetDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Modal Demos";
        self.buttonLink = @"https://docs.criticalmoments.io/actions/modals";
        self.infoText = @"Announcements, decisions, and more. Powerful and beautiful messaging with native modal UI";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {

    // Basics

    CMDemoAction *announceSheet = [[CMDemoAction alloc] init];
    announceSheet.title = @"Announcement";
    announceSheet.subtitle = @"Display a sheet with an announcement for the user, including a custom action button.";
    announceSheet.actionCMActionName = @"simpleModalAction";
    [announceSheet addResetTestTarget:self action:@selector(dismissSheets)];

    CMDemoAction *tip = [[CMDemoAction alloc] init];
    tip.title = @"Tip";
    tip.subtitle = @"Show a sample tip for the user.";
    tip.actionCMActionName = @"headphoneModalExample";
    [tip addResetTestTarget:self action:@selector(dismissSheets)];

    CMDemoAction *theme = [[CMDemoAction alloc] init];
    theme.title = @"Themed Example";
    theme.subtitle = @"Show an announcement with a custom theme.";
    theme.actionCMActionName = @"themeModalExample";
    [theme addResetTestTarget:self action:@selector(dismissSheets)];

    [self addSection:@"Examples" withActions:@[ announceSheet, tip, theme ]];

    CMDemoAction *imageSheet = [[CMDemoAction alloc] init];
    imageSheet.title = @"Image style options";
    imageSheet.subtitle = @"Display a variety of image styles";
    imageSheet.actionCMActionName = @"imageStylesModal";
    [imageSheet addResetTestTarget:self action:@selector(dismissSheets)];

    CMDemoAction *typography = [[CMDemoAction alloc] init];
    typography.title = @"Typography options";
    typography.subtitle = @"Show a variety of typography.";
    typography.actionCMActionName = @"typographyModalExample";
    [typography addResetTestTarget:self action:@selector(dismissSheets)];

    CMDemoAction *buttons = [[CMDemoAction alloc] init];
    buttons.title = @"Button options";
    buttons.subtitle = @"Show the buton styling options.";
    buttons.actionCMActionName = @"buttonsModalExample";
    [buttons addResetTestTarget:self action:@selector(dismissSheets)];

    [self addSection:@"Style Options" withActions:@[ imageSheet, typography, buttons ]];
}

- (void)dismissSheets {
    [Utils.keyWindow.rootViewController.presentedViewController dismissViewControllerAnimated:NO completion:nil];
}

@end
