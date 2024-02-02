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
        self.infoText = @"Asking users to subscribe or buy before "
                        @"they have experienced the core value of your app is a sure fire way to get them to decline, "
                        @"or even worse, leave for good.\n\nCritical Moments helps you create a user journey:\n1) "
                        @"Ensure they discover and use the core features of your app\n2) After they have seen value, "
                        @"ask them to subscribe or buy.\n3) Once they are loyal, ask them to review.\n\nThe example "
                        @"below walk through the user journey of a fictional “todo "
                        @"list” app.";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {
}

@end
