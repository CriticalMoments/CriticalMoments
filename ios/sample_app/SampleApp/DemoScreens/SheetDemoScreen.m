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
        self.title = @"Sheet Demos";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {

    // Basics

    CMDemoAction *announceSheet = [[CMDemoAction alloc] init];
    announceSheet.title = @"Annoucement";
    announceSheet.subtitle = @"Display a sheet with an announcement for the user.";
    announceSheet.actionCMActionName = @"simpleModalAction";
    [announceSheet addResetTestTarget:self action:@selector(dismissSheets)];

    CMDemoAction *tip = [[CMDemoAction alloc] init];
    tip.title = @"Tip";
    tip.subtitle = @"Show a sample tip for the user.";
    tip.actionCMActionName = @"headphoneModalExample";
    [tip addResetTestTarget:self action:@selector(dismissSheets)];

    [self addSection:@"Examples" withActions:@[ announceSheet, tip ]];

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

    [self addSection:@"Style Options" withActions:@[ imageSheet, typography ]];
}

- (void)dismissSheets {
    [Utils.keyWindow.rootViewController.presentedViewController dismissViewControllerAnimated:NO completion:nil];
}

@end
