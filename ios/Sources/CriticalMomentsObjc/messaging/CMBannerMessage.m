//
//  CMBannerMessage.m
//  
//
//  Created by Steve Cosman on 2023-04-22.
//

#import "CMBannerMessage.h"

@import UIKit;

@interface CMBannerMessage ()

@property (nonatomic, strong, readwrite) NSString* body;

@end

@implementation CMBannerMessage

-(instancetype)initWithBody:(NSString*)body {
    self = [super init];
    if (self) {
        self.body = body;
        self.showDismissButton = YES;
    }
    return self;
}


-(UIView*) buildViewForMessage {
    UIView* view = [[UIView alloc] init];
    
    // TODO: load from theme
    UIColor* forgroundBannerColor = [UIColor blackColor];
    UIColor* backgroundBannerColor = [UIColor systemYellowColor];
    
    view.backgroundColor = backgroundBannerColor;
    
    UILabel* bodyLabel = [[UILabel alloc] init];
    bodyLabel.text = self.body;
    bodyLabel.textColor = forgroundBannerColor;
    bodyLabel.backgroundColor = [UIColor clearColor];
    bodyLabel.translatesAutoresizingMaskIntoConstraints = NO;
    bodyLabel.numberOfLines = 2;
    // TODO style
    // TODO elipisis
    // TODO height passed up
    [view addSubview:bodyLabel];
    
    // Layout
    NSArray<NSLayoutConstraint*>* constraints = @[
        [bodyLabel.topAnchor constraintEqualToAnchor:view.layoutMarginsGuide.topAnchor],
        [bodyLabel.leftAnchor constraintEqualToAnchor:view.layoutMarginsGuide.leftAnchor],
        [bodyLabel.bottomAnchor constraintEqualToAnchor:view.layoutMarginsGuide.bottomAnchor],
    ];
    [NSLayoutConstraint activateConstraints:constraints];
    
    if (!self.showDismissButton) {
        // Layout body without dismiss button
        constraints = [constraints arrayByAddingObject:[bodyLabel.rightAnchor constraintEqualToAnchor:view.layoutMarginsGuide.rightAnchor]];
    } else {
        UIButton* dismissButton = [UIButton buttonWithType:UIButtonTypeCustom];
        if (@available(iOS 13.0, *)) {
            UIImage *dismissImage = [[UIImage systemImageNamed:@"xmark"] imageWithTintColor:forgroundBannerColor renderingMode:UIImageRenderingModeAlwaysOriginal];
            [dismissButton setImage:dismissImage forState:UIControlStateNormal];
        } else {
            [dismissButton setTitle:@"X" forState:UIControlStateNormal];
            [dismissButton setTitleColor:forgroundBannerColor forState:UIControlStateNormal];
        }
        [dismissButton addTarget:self action:@selector(dismissTapped:) forControlEvents:UIControlEventPrimaryActionTriggered];
        dismissButton.translatesAutoresizingMaskIntoConstraints = NO;
        [view addSubview:dismissButton];

        // Layout for dismiss button, making room for body. 44=HIG accessibility recommendation.
        constraints = [constraints arrayByAddingObjectsFromArray:@[
            [dismissButton.heightAnchor constraintEqualToConstant:44],
            [dismissButton.widthAnchor constraintEqualToConstant:44],
            [dismissButton.rightAnchor constraintEqualToAnchor:view.layoutMarginsGuide.rightAnchor],
            [dismissButton.centerYAnchor constraintEqualToAnchor:view.layoutMarginsGuide.centerYAnchor],
            [bodyLabel.rightAnchor constraintEqualToAnchor:dismissButton.leftAnchor constant:-12],
        ]];
    }
    
    [NSLayoutConstraint activateConstraints:constraints];
    
    return view;
}

- (void)dismissTapped:(UIButton*)sender {
    [self.dismissDelegate dismissedMessage:self];
}

@end
