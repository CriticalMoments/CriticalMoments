//
//  MainDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import "MainDemoScreen.h"

#import "BannerDemoScreen.h"

@implementation MainDemoScreen

-(instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Critical Moments";
        [self buildSections];
    }
    return self;
}

-(void) buildSections {
    CMDemoAction* bannersAction = [[CMDemoAction alloc] init];
    bannersAction.title = @"Banners";
    bannersAction.subtitle = @"UI to display announcement banners across the top or bottom of your app.";
    bannersAction.actionNextScreen = [[BannerDemoScreen alloc] init];
    
    [self addSection:@"UI Actions" withActions:@[
        bannersAction
    ]];
}

@end
