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

    [self addSection:@"Web links"
         withActions:@[
             openWebLink,
         ]];
}

@end
