//
//  LinkDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-05-11.
//

#import "LinkDemoScreen.h"

@implementation LinkDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Link Demos";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {

    // Web links

    CMDemoAction *openWebLink = [[CMDemoAction alloc] init];
    openWebLink.title = @"Open web link";
    openWebLink.subtitle = @"Open criticalmoments.io in your web browser";
    openWebLink.actionCMActionName = @"web_link_action";

    // TODO another web link, local browser VC

    [self addSection:@"Web links"
         withActions:@[
             openWebLink,
         ]];

    // TODO: a app deeplink

    // App deep links

    CMDemoAction *openSettingsLink = [[CMDemoAction alloc] init];
    openSettingsLink.title = @"Open settings";
    openSettingsLink.subtitle =
        @"Open an app deeplink into the iOS Settings app";
    openSettingsLink.actionCMActionName = @"settings_link_action";

    CMDemoAction *openMainScreenLink = [[CMDemoAction alloc] init];
    openMainScreenLink.title = @"Open deeplink";
    openMainScreenLink.subtitle =
        @"Open an app deeplink into this sample app's main screen";
    openMainScreenLink.actionCMActionName = @"main_screen_deeplink_action";

    [self addSection:@"App deep links"
         withActions:@[ openSettingsLink, openMainScreenLink ]];
}

@end
