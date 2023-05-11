//
//  MainDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import "MainDemoScreen.h"

#import "AlertDemoScreen.h"
#import "BannerDemoScreen.h"
#import "ThemeDemoScreen.h"

@implementation MainDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Critical Moments";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {
    CMDemoAction *bannersAction = [[CMDemoAction alloc] init];
    bannersAction.title = @"Banners";
    bannersAction.subtitle = @"UI to display announcement banners across the "
                             @"top or bottom of your app.";
    bannersAction.actionNextScreen = [[BannerDemoScreen alloc] init];

    CMDemoAction *alertAction = [[CMDemoAction alloc] init];
    alertAction.title = @"Alerts";
    alertAction.subtitle = @"UI to display system alerts/dialogs.";
    alertAction.actionNextScreen = [[AlertDemoScreen alloc] init];

    [self addSection:@"UI Actions" withActions:@[ bannersAction, alertAction ]];

    CMDemoAction *themeAction = [[CMDemoAction alloc] init];
    themeAction.title = @"Edit Theme";
    themeAction.subtitle = @"Modify the colors, font and style of UI elements.";
    themeAction.actionNextScreen = [[ThemeDemoScreen alloc] init];

    [self addSection:@"Themes / Style" withActions:@[ themeAction ]];
}

@end
