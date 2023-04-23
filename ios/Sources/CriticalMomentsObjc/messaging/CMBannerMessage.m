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
        [self buildSubviews];
    }
    return self;
}

-(void)setDismissButton:(UIButton *)dismissButton {
    [dismissButton addTarget:self action:@selector(dismissTapped:) forControlEvents:UIControlEventPrimaryActionTriggered];
    _dismissButton = dismissButton;
}

-(void) buildSubviews {
    // TODO
    self.backgroundColor = [UIColor greenColor];
    
    self.bodyLabel = [[UILabel alloc] init];
    self.bodyLabel.text = self.body;
    self.bodyLabel.backgroundColor = [UIColor clearColor];
    self.bodyLabel.translatesAutoresizingMaskIntoConstraints = NO;
    self.bodyLabel.numberOfLines = 2;
    // TODO style
    // TODO elipisis
    // TODO height passed up
    [self addSubview:self.bodyLabel];
    
    // TODO Warning
    // TODO style/color
    self.dismissButton = [UIButton buttonWithType:UIButtonTypeClose];
    self.dismissButton.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:self.dismissButton];
    
    [self setupLayout];
}

-(void) setupLayout {
    // Layout
    NSArray<NSLayoutConstraint*>* constraints = @[
        [self.dismissButton.heightAnchor constraintEqualToConstant:40],
        [self.dismissButton.widthAnchor constraintEqualToConstant:40],
        [self.dismissButton.rightAnchor constraintEqualToAnchor:self.layoutMarginsGuide.rightAnchor],
        [self.dismissButton.centerYAnchor constraintEqualToAnchor:self.layoutMarginsGuide.centerYAnchor],
        [self.bodyLabel.topAnchor constraintEqualToAnchor:self.layoutMarginsGuide.topAnchor],
        [self.bodyLabel.leftAnchor constraintEqualToAnchor:self.layoutMarginsGuide.leftAnchor],
        [self.bodyLabel.rightAnchor constraintEqualToAnchor:self.dismissButton.leftAnchor constant:-12],
        [self.bodyLabel.bottomAnchor constraintEqualToAnchor:self.layoutMarginsGuide.bottomAnchor],
    ];
    [NSLayoutConstraint activateConstraints:constraints];
}

- (void)dismissTapped:(UIButton*)sender {
    [self.delegate dismissed];
}

@end
