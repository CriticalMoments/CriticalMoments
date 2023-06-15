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
    announceSheet.title = @"Show Annoucement Sheet";
    announceSheet.subtitle = @"Display a sheet with announement/information for the user";
    // announceSheet.actionCMActionName = @"show_notice_alert";
    announceSheet.actionBlock = ^{
      CMSheetViewController *sheetVc = [[CMSheetViewController alloc] init];
      [Utils.keyWindow.rootViewController presentViewController:sheetVc animated:YES completion:nil];
    };
    [announceSheet addResetTestTarget:self action:@selector(dismissSheets)];
    [self addActionToRootSection:announceSheet];
}

- (void)dismissSheets {
    [Utils.keyWindow.rootViewController.presentedViewController dismissViewControllerAnimated:YES completion:nil];
}

@end
