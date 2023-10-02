//
//  CMBannerMessage.m
//
//
//  Created by Steve Cosman on 2023-04-22.
//

#import "CMBannerMessage.h"
#import "../include/CriticalMoments.h"
#import "../themes/CMTheme.h"
#import "../themes/CMTheme_private.h"
#import "CMBannerMessage_private.h"

@import UIKit;

@interface CMBannerMessage () <CMBannerActionDelegate>

@property(nonatomic, strong, readwrite) NSString *body;
@property(nonatomic, strong, readwrite) NSString *appcoreTapActionName;
@property(nonatomic, strong, readwrite) NSString *appcoreThemeName;
@property(nonnull, strong, readwrite) UIButton *nextButton, *dismissButton;
@property(nonnull, strong, readwrite) UILabel *bodyLabel;

@end

@implementation CMBannerMessage

- (instancetype)initWithBody:(NSString *)body {
    self = [super init];
    if (self) {
        self.body = body;
        self.showDismissButton = YES;
        self.preferredPosition = CMBannerPositionNoPreference;

        [self buildView];
    }
    return self;
}

- (instancetype)initWithAppcoreDataModel:(DatamodelBannerAction *)bannerData {
    self = [super init];
    if (self) {
        self.body = bannerData.body;

        self.showDismissButton = bannerData.showDismissButton;
        if (bannerData.maxLineCount != DatamodelBannerMaxLineCountSystemDefault) {
            self.maxLineCount = @(bannerData.maxLineCount);
        }
        if (bannerData.tapActionName.length > 0) {
            self.appcoreTapActionName = bannerData.tapActionName;
            self.actionDelegate = self;
        }
        if (bannerData.customThemeName.length > 0) {
            CMTheme *customTheme = [CMTheme namedThemeFromAppcore:bannerData.customThemeName];
            self.customTheme = customTheme;
        }

        self.preferredPosition = CMBannerPositionNoPreference;
        if ([DatamodelBannerPositionTop isEqualToString:bannerData.preferredPosition]) {
            self.preferredPosition = CMBannerPositionTop;
        } else if ([DatamodelBannerPositionBottom isEqualToString:bannerData.preferredPosition]) {
            self.preferredPosition = CMBannerPositionBottom;
        }

        [self buildView];
    }
    return self;
}

- (void)layoutSubviews {
    [super layoutSubviews];
    [self.messageManagerDelegate messageDidLayout:self];
}

- (void)buildView {
    CMTheme *theme = self.customTheme;
    if (!theme) {
        theme = CMTheme.current;
    }
    UIColor *forgroundBannerColor = theme.bannerForegroundColor;
    UIColor *backgroundBannerColor = theme.bannerBackgroundColor;
    UIFont *bannerFont = [theme boldFontOfSize:theme.bodyFontSize];

    self.backgroundColor = backgroundBannerColor;

    _bodyLabel = [[UILabel alloc] init];
    _bodyLabel.text = self.body;
    _bodyLabel.textColor = forgroundBannerColor;
    _bodyLabel.font = bannerFont;
    _bodyLabel.backgroundColor = [UIColor clearColor];
    _bodyLabel.translatesAutoresizingMaskIntoConstraints = NO;
    if (self.maxLineCount) {
        _bodyLabel.numberOfLines = [self.maxLineCount intValue];
    } else {
        _bodyLabel.numberOfLines = 4;
    }
    [self addSubview:_bodyLabel];

    // Gesture for action
    UITapGestureRecognizer *tapReco = [[UITapGestureRecognizer alloc] initWithTarget:self
                                                                              action:@selector(bannerTapped)];
    [self setUserInteractionEnabled:YES];
    [self addGestureRecognizer:tapReco];

    // Layout

    NSArray<NSLayoutConstraint *> *constraints = @[
        [_bodyLabel.topAnchor constraintEqualToAnchor:self.layoutMarginsGuide.topAnchor],
        [_bodyLabel.bottomAnchor constraintEqualToAnchor:self.layoutMarginsGuide.bottomAnchor],
        [_bodyLabel.leftAnchor constraintGreaterThanOrEqualToAnchor:self.layoutMarginsGuide.leftAnchor],
        [_bodyLabel.rightAnchor constraintLessThanOrEqualToAnchor:self.layoutMarginsGuide.rightAnchor],
        [_bodyLabel.centerXAnchor constraintEqualToAnchor:self.centerXAnchor],
        [_bodyLabel.widthAnchor constraintLessThanOrEqualToConstant:CM_MAX_TEXT_WIDTH],
    ];
    [NSLayoutConstraint activateConstraints:constraints];

    _dismissButton = [UIButton buttonWithType:UIButtonTypeCustom];
    if (@available(iOS 13.0, *)) {
        UIImage *dismissImage =
            [[UIImage systemImageNamed:@"xmark"] imageWithTintColor:forgroundBannerColor
                                                      renderingMode:UIImageRenderingModeAlwaysOriginal];
        [_dismissButton setImage:dismissImage forState:UIControlStateNormal];
    } else {
        [_dismissButton setTitle:@"✕" forState:UIControlStateNormal];
        [_dismissButton setTitleColor:forgroundBannerColor forState:UIControlStateNormal];
    }
    [_dismissButton addTarget:self
                       action:@selector(dismissTapped:)
             forControlEvents:UIControlEventPrimaryActionTriggered];
    _dismissButton.translatesAutoresizingMaskIntoConstraints = NO;
    [self updateDismissButonState];

    // Create "ᐊᐅ" button for next message
    _nextButton = [UIButton buttonWithType:UIButtonTypeCustom];
    [_nextButton setTitle:@"ᐊᐅ" forState:UIControlStateNormal];
    [_nextButton setTitleColor:forgroundBannerColor forState:UIControlStateNormal];
    [_nextButton addTarget:self
                    action:@selector(nextMessageButtonTapped:)
          forControlEvents:UIControlEventPrimaryActionTriggered];
    _nextButton.translatesAutoresizingMaskIntoConstraints = NO;
    [self updateNextButtonState];
}

- (void)setShowDismissButton:(bool)showDismissButton {
    _showDismissButton = showDismissButton;
    [self updateDismissButonState];
}

- (void)setNextMessageDelegate:(id<CMBannerNextMessageDelegate>)nextMessageDelegate {
    _nextMessageDelegate = nextMessageDelegate;
    [self updateNextButtonState];
}

- (void)updateDismissButonState {
    if (!_dismissButton) {
        return;
    }

    if (!self.showDismissButton && _dismissButton.superview) {
        [_dismissButton removeFromSuperview];
    } else if (self.showDismissButton && _dismissButton.superview == nil) {
        // Layout for dismiss button, making room for body
        [self addSubview:_dismissButton];

        NSArray<NSLayoutConstraint *> *constraints = @[
            // 44=HIG accessibility recommendation
            [_dismissButton.heightAnchor constraintEqualToConstant:44],
            [_dismissButton.widthAnchor constraintEqualToConstant:44],
            [_dismissButton.rightAnchor constraintEqualToAnchor:self.safeAreaLayoutGuide.rightAnchor],
            [_dismissButton.centerYAnchor constraintEqualToAnchor:self.safeAreaLayoutGuide.centerYAnchor],
            // -6 is just visual padding on left of X button
            [_bodyLabel.rightAnchor constraintLessThanOrEqualToAnchor:_dismissButton.leftAnchor constant:-6],
        ];
        [NSLayoutConstraint activateConstraints:constraints];
    }
}

- (void)updateNextButtonState {
    if (!_nextButton) {
        return;
    }

    if (_nextMessageDelegate == nil && _nextButton.superview != nil) {
        [_nextButton removeFromSuperview];
    } else if (_nextMessageDelegate != nil && _nextButton.superview == nil) {
        [self addSubview:_nextButton];

        NSArray<NSLayoutConstraint *> *constraints = @[
            [_nextButton.heightAnchor constraintGreaterThanOrEqualToConstant:44],
            [_nextButton.widthAnchor constraintGreaterThanOrEqualToConstant:44],
            [_nextButton.leftAnchor constraintEqualToAnchor:self.safeAreaLayoutGuide.leftAnchor],
            [_nextButton.centerYAnchor constraintEqualToAnchor:self.safeAreaLayoutGuide.centerYAnchor],
            [_bodyLabel.leftAnchor constraintGreaterThanOrEqualToAnchor:_nextButton.rightAnchor],
        ];
        [NSLayoutConstraint activateConstraints:constraints];
    }
}

- (void)dismissTapped:(UIButton *)sender {
    [self.messageManagerDelegate dismissedMessage:self];
}

- (void)nextMessageButtonTapped:(UIButton *)sender {
    [self.nextMessageDelegate nextMessage];
}

- (void)bannerTapped {
    [self.actionDelegate messageAction:self];
}

#pragma mark CMBannerActionDelegate

// Messages should only be their own delegate for DatamodelBannerAction messages
- (void)messageAction:(nonnull CMBannerMessage *)message {
    if (!self.appcoreTapActionName) {
        return;
    }

    NSError *error;
    [CriticalMoments.sharedInstance performNamedAction:self.appcoreTapActionName error:&error];
    if (error) {
        NSLog(@"CriticalMoments: Banner tap unknown issue: %@", error);
    }
}

@end
