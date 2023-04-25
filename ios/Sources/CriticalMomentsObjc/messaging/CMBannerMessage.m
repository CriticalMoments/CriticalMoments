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
@property (nonatomic, strong) UILabel* bodyLabel;
@property (nonatomic, strong) UIButton* dismissButton;

@end

@implementation CMBannerMessage

-(instancetype)initWithBody:(NSString*)body {
    self = [super init];
    if (self) {
        self.body = body;
    }
    return self;
}

-(void)setDismissButton:(UIButton *)dismissButton {
    [dismissButton addTarget:self action:@selector(dismissTapped:) forControlEvents:UIControlEventPrimaryActionTriggered];
    _dismissButton = dismissButton;
}

-(UIView*) buildViewForMessage {
    UIView* view = [[UIView alloc] init];
    
    // TODO: load from theme
    UIColor* forgroundBannerColor = [UIColor blackColor];
    UIColor* backgroundBannerColor = [UIColor greenColor];
    
    view.backgroundColor = backgroundBannerColor;
    
    self.bodyLabel = [[UILabel alloc] init];
    self.bodyLabel.text = self.body;
    self.bodyLabel.textColor = forgroundBannerColor;
    self.bodyLabel.backgroundColor = [UIColor clearColor];
    self.bodyLabel.translatesAutoresizingMaskIntoConstraints = NO;
    self.bodyLabel.numberOfLines = 2;
    // TODO style
    // TODO elipisis
    // TODO height passed up
    [view addSubview:self.bodyLabel];
    
    // TODO style/color
    self.dismissButton = [UIButton buttonWithType:UIButtonTypeCustom];
    if (@available(iOS 13.0, *)) {
        UIImage *dismissImage = [[UIImage systemImageNamed:@"xmark"] imageWithTintColor:forgroundBannerColor renderingMode:UIImageRenderingModeAlwaysOriginal];
        [self.dismissButton setImage:dismissImage forState:UIControlStateNormal];
    } else {
        // Fallback on earlier versions
        [self.dismissButton setTitle:@"X" forState:UIControlStateNormal];
        [self.dismissButton setTitleColor:forgroundBannerColor forState:UIControlStateNormal];
    }
    
    self.dismissButton.translatesAutoresizingMaskIntoConstraints = NO;
    [view addSubview:self.dismissButton];
    
    [self setupLayoutForRootView:(UIView*)view];
    
    return view;
}

-(void) setupLayoutForRootView:(UIView*)view {
    // Layout
    NSArray<NSLayoutConstraint*>* constraints = @[
        [self.dismissButton.heightAnchor constraintEqualToConstant:40],
        [self.dismissButton.widthAnchor constraintEqualToConstant:40],
        [self.dismissButton.rightAnchor constraintEqualToAnchor:view.layoutMarginsGuide.rightAnchor],
        [self.dismissButton.centerYAnchor constraintEqualToAnchor:view.layoutMarginsGuide.centerYAnchor],
        [self.bodyLabel.topAnchor constraintEqualToAnchor:view.layoutMarginsGuide.topAnchor],
        [self.bodyLabel.leftAnchor constraintEqualToAnchor:view.layoutMarginsGuide.leftAnchor],
        [self.bodyLabel.rightAnchor constraintEqualToAnchor:self.dismissButton.leftAnchor constant:-12],
        [self.bodyLabel.bottomAnchor constraintEqualToAnchor:view.layoutMarginsGuide.bottomAnchor],
    ];
    [NSLayoutConstraint activateConstraints:constraints];
}

- (void)dismissTapped:(UIButton*)sender {
    [self.dismissDelegate dismissedMessage:self];
}

@end
