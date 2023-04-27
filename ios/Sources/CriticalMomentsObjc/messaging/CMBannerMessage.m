//
//  CMBannerMessage.m
//  
//
//  Created by Steve Cosman on 2023-04-22.
//

#import "CMBannerMessage.h"
#import "CMBannerMessage_private.h"

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
    UIFont* bannerFont = [UIFont boldSystemFontOfSize:[UIFont systemFontSize]];
    
    view.backgroundColor = backgroundBannerColor;
    
    UILabel* bodyLabel = [[UILabel alloc] init];
    bodyLabel.text = self.body;
    bodyLabel.textColor = forgroundBannerColor;
    bodyLabel.font = bannerFont;
    bodyLabel.backgroundColor = [UIColor clearColor];
    bodyLabel.translatesAutoresizingMaskIntoConstraints = NO;
    if (self.maxLineCount) {
        bodyLabel.numberOfLines = [self.maxLineCount intValue];
    } else {
        bodyLabel.numberOfLines = 4;
    }
    [view addSubview:bodyLabel];
    
    // Gesture for action
    if (self.actionDelegate) {
        UITapGestureRecognizer* tapReco = [[UITapGestureRecognizer alloc] initWithTarget:self action:@selector(bannerTapped)];
        [view setUserInteractionEnabled:YES];
        [view addGestureRecognizer:tapReco];
    }
    
    // Layout
    
    NSArray<NSLayoutConstraint*>* constraints = @[
        [bodyLabel.topAnchor constraintEqualToAnchor:view.layoutMarginsGuide.topAnchor],
        [bodyLabel.bottomAnchor constraintEqualToAnchor:view.layoutMarginsGuide.bottomAnchor],
        [bodyLabel.leftAnchor constraintGreaterThanOrEqualToAnchor:view.layoutMarginsGuide.leftAnchor],
        [bodyLabel.rightAnchor constraintLessThanOrEqualToAnchor:view.layoutMarginsGuide.rightAnchor],
        [bodyLabel.centerXAnchor constraintEqualToAnchor:view.centerXAnchor],
        // Max width for iPad, based on readableContentGuide from Apple
        [bodyLabel.widthAnchor constraintLessThanOrEqualToConstant:672],
    ];
    [NSLayoutConstraint activateConstraints:constraints];
    
    if (self.showDismissButton) {
        UIButton* dismissButton = [UIButton buttonWithType:UIButtonTypeCustom];
        if (@available(iOS 13.0, *)) {
            UIImage *dismissImage = [[UIImage systemImageNamed:@"xmark"] imageWithTintColor:forgroundBannerColor renderingMode:UIImageRenderingModeAlwaysOriginal];
            [dismissButton setImage:dismissImage forState:UIControlStateNormal];
        } else {
            [dismissButton setTitle:@"✕" forState:UIControlStateNormal];
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
            [bodyLabel.rightAnchor constraintLessThanOrEqualToAnchor:dismissButton.leftAnchor constant:-4],
        ]];
    }
    
    // Check for multiple messages
    if (self.nextMessageDelegate) {
        // Create "ᐊᐅ" button
        UIButton* nextMessageButton = [UIButton buttonWithType:UIButtonTypeCustom];
        [nextMessageButton setTitle:@"ᐊᐅ" forState:UIControlStateNormal];
        [nextMessageButton setTitleColor:forgroundBannerColor forState:UIControlStateNormal];
        [nextMessageButton addTarget:self action:@selector(nextTapped:) forControlEvents:UIControlEventPrimaryActionTriggered];
        nextMessageButton.translatesAutoresizingMaskIntoConstraints = NO;
        [view addSubview:nextMessageButton];
        
        // 44=HIG min size for tap target
        constraints = [constraints arrayByAddingObjectsFromArray:@[
            [nextMessageButton.heightAnchor constraintGreaterThanOrEqualToConstant:44],
            [nextMessageButton.widthAnchor constraintGreaterThanOrEqualToConstant:44],
            [nextMessageButton.leftAnchor constraintEqualToAnchor:view.layoutMarginsGuide.leftAnchor],
            [nextMessageButton.centerYAnchor constraintEqualToAnchor:view.layoutMarginsGuide.centerYAnchor],
            [bodyLabel.leftAnchor constraintGreaterThanOrEqualToAnchor:nextMessageButton.rightAnchor ],
        ]];
    }
    
    [NSLayoutConstraint activateConstraints:constraints];
    
    return view;
}

- (void)dismissTapped:(UIButton*)sender {
    [self.dismissDelegate dismissedMessage:self];
}

- (void)nextTapped:(UIButton*)sender {
    [self.nextMessageDelegate nextMessage];
}

- (void)bannerTapped {
    [self.actionDelegate messageAction:self];
}

@end
