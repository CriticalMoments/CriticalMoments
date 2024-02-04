//
//  MainDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import "MainDemoScreen.h"

#import "AlertDemoScreen.h"
#import "BannerDemoScreen.h"
#import "ConditionsDemoScreen.h"
#import "ConversionDemoScreen.h"
#import "LinkDemoScreen.h"
#import "SheetDemoScreen.h"
#import "ThemeDemoScreen.h"

@implementation MainDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Critical Moments";
        self.infoText = @"Explore demos of the Critial Moments SDK";
        self.buttonLink = @"https://docs.criticalmoments.io";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {
    CMDemoAction *conversionDemo = [[CMDemoAction alloc] init];
    conversionDemo.title = @"Conversion & User Journey";
    conversionDemo.subtitle = @"Help your users get the most from your app, and convert to paid.";
    conversionDemo.actionNextScreen = [[ConversionDemoScreen alloc] init];

    [self addSection:@"Use Case Examples" withActions:@[ conversionDemo ]];

    CMDemoAction *conditionDemos = [[CMDemoAction alloc] init];
    conditionDemos.title = @"Conditions";
    conditionDemos.subtitle = @"Evaluate powerful conditions at runtime";
    conditionDemos.actionNextScreen = [[ConditionsDemoScreen alloc] init];
    conditionDemos.skipInUiTesting = true;

    [self addSection:@"Conditions" withActions:@[ conditionDemos ]];

    CMDemoAction *bannersAction = [[CMDemoAction alloc] init];
    bannersAction.title = @"Banners";
    bannersAction.subtitle = @"UI to display announcement banners across the "
                             @"top or bottom of your app.";
    bannersAction.actionNextScreen = [[BannerDemoScreen alloc] init];

    CMDemoAction *alertAction = [[CMDemoAction alloc] init];
    alertAction.title = @"Alerts";
    alertAction.subtitle = @"UI to display system alerts and action sheets.";
    alertAction.actionNextScreen = [[AlertDemoScreen alloc] init];

    CMDemoAction *linksAction = [[CMDemoAction alloc] init];
    linksAction.title = @"Links";
    linksAction.subtitle = @"Open web links or app deeplinks";
    linksAction.actionNextScreen = [[LinkDemoScreen alloc] init];

    CMDemoAction *sheetAction = [[CMDemoAction alloc] init];
    sheetAction.title = @"Modals";
    sheetAction.subtitle = @"UI to display announcements or other information in sheets which overlay the entire app.";
    sheetAction.actionNextScreen = [[SheetDemoScreen alloc] init];

    [self addSection:@"Actions" withActions:@[ bannersAction, alertAction, linksAction, sheetAction ]];

    CMDemoAction *themeAction = [[CMDemoAction alloc] init];
    themeAction.title = @"Edit Theme";
    themeAction.subtitle = @"Modify the colors, font and style of UI elements.";
    themeAction.actionNextScreen = [[ThemeDemoScreen alloc] init];

    [self addSection:@"Themes / Style" withActions:@[ themeAction ]];
}

@end
