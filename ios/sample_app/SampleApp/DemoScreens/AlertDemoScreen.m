//
//  AlertDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-05-11.
//

#import "AlertDemoScreen.h"

@import CriticalMoments;

@implementation AlertDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Alerts Demos";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {

    // Basics

    CMDemoAction *noticeAlert = [[CMDemoAction alloc] init];
    noticeAlert.title = @"Show Notice Alert";
    noticeAlert.subtitle = @"Display a short alert with OK button";
    noticeAlert.actionCMActionName = @"show_notice_alert";
    [self addActionToRootSection:noticeAlert];

    CMDemoAction *cancelAlert = [[CMDemoAction alloc] init];
    cancelAlert.title = @"Show Cancelable Alert";
    cancelAlert.subtitle = @"Display a short alert with OK and Cancel buttons";
    cancelAlert.actionCMActionName = @"show_cancelable_alert";
    [self addActionToRootSection:cancelAlert];

    CMDemoAction *largeAlert = [[CMDemoAction alloc] init];
    largeAlert.title = @"Show Action Sheet Alert";
    largeAlert.subtitle =
        @"Display a sheet style alert, with custom buttons and actions";
    largeAlert.actionCMActionName = @"custom_button_alert_large";
    [self addActionToRootSection:largeAlert];

    CMDemoAction *severalButtonAlert = [[CMDemoAction alloc] init];
    severalButtonAlert.title = @"Show many button alert";
    severalButtonAlert.subtitle =
        @"Display a alert, with custom buttons and actions";
    severalButtonAlert.actionCMActionName = @"custom_button_alert_dialog";
    [self addActionToRootSection:severalButtonAlert];
}

@end
