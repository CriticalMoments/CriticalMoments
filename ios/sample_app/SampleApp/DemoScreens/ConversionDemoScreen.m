//
//  ConversionDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2024-02-02.
//

#import "ConversionDemoScreen.h"

@implementation ConversionDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Conversions & Journey";
        self.infoText = @"It's important to help users discover the amazing features your app has to offer. This makes "
                        @"them much more likely to become long term users (activation, retention), and satisfied users "
                        @"(subscribe/buy, higher ratings, lower churn).\n\nAsking users to subscribe or buy before "
                        @"they have experienced the core value of your app, is a sure fire way to get them to decline "
                        @"or even leave for good.\n\nThese examples walk through the user journey of a fictional “todo "
                        @"list” app. In the early steps, we nudge our the user to discover the app and its key "
                        @"features/value-prop. Once they have we attempt to get them to subscribe and review.";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {
}

@end
