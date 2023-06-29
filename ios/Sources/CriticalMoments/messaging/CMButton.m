//
//  CMButtonStack.m
//
//
//  Created by Steve Cosman on 2023-06-15.
//

#import "CMButton.h"

// Align to default and size=small OS buttons
#define CM_OS_BUTTON_FONT_SIZE 17.0
#define CM_SMALL_BUTTON_FONT_SIZE 15.0

@interface CMButton ()

@property(nonatomic, strong) CMTheme *customTheme;
@property(nonatomic, strong) DatamodelButton *model;
@property(nonatomic, strong) UIButton *buttonView;

@end

@implementation CMButton

- (instancetype)initWithWithDataModel:(DatamodelButton *)model andTheme:(CMTheme *)theme {
    self = [super init];
    if (self) {
        _model = model;
        _customTheme = theme;

        self.buttonView = [self buttonWithWithDataModel:self.model];
        [self addSubview:self.buttonView];

        [self.buttonView addTarget:self
                            action:@selector(buttonTapped:)
                  forControlEvents:UIControlEventPrimaryActionTriggered];
    }
    return self;
}

- (CMTheme *)theme {
    if (self.customTheme) {
        return self.customTheme;
    }
    return CMTheme.current;
}

- (void)layoutSubviews {
    [super layoutSubviews];
    self.buttonView.frame = self.bounds;
}

- (CGSize)intrinsicContentSize {
    return self.buttonView.intrinsicContentSize;
}

- (void)buttonTapped:(UIButton *)target {
    // two actions fire:
    // 1) the model's named action, if it exists
    // 2) the system's default action, unless prevent default (eg: dismiss modal)

    if (self.model.actionName.length > 0) {
        NSError *error;
        [AppcoreSharedAppcore() performNamedAction:self.model.actionName error:&error];
        if (error) {
            NSLog(@"CriticalMoments: Button tap unknown issue: %@", error);
        }
    }

    if (!self.model.preventDefault && self.defaultAction) {
        self.defaultAction();
    }
}

- (UIButton *)buttonWithWithDataModel:(DatamodelButton *)model {
    UIButton *button = [self buildNewStyleButtonWithDataModel:model];
    if (!button) {
        // earlier versions of iOS
        button = [self buildOldStyleButtonWithDataModel:model];
    }

    [button setTitle:model.title forState:UIControlStateNormal];

    return button;
}

- (UIButton *)buildNewStyleButtonWithDataModel:(DatamodelButton *)model {
    if (@available(iOS 15.0, *)) {
        UIButtonConfiguration *c;
        if ([DatamodelButtonStyleEnumLarge isEqualToString:model.style]) {
            c = UIButtonConfiguration.filledButtonConfiguration;
            c.buttonSize = UIButtonConfigurationSizeLarge;
            c.baseBackgroundColor = [self.theme primaryColorForView:self];
        } else if ([DatamodelButtonStyleEnumSecondary isEqualToString:model.style]) {
            c = UIButtonConfiguration.tintedButtonConfiguration;
            c.baseBackgroundColor = [self.theme primaryColorForView:self];
            c.baseForegroundColor = [self.theme primaryColorForView:self];
        } else if ([DatamodelButtonStyleEnumTertiary isEqualToString:model.style]) {
            c = UIButtonConfiguration.grayButtonConfiguration;
            c.baseForegroundColor = [self.theme primaryColorForView:self];
        } else if ([DatamodelButtonStyleEnumInfo isEqualToString:model.style]) {
            c = UIButtonConfiguration.plainButtonConfiguration;
            c.baseForegroundColor = [self.theme primaryColorForView:self];
        } else if ([DatamodelButtonStyleEnumInfoSmall isEqualToString:model.style]) {
            c = UIButtonConfiguration.plainButtonConfiguration;
            c.baseForegroundColor = [self.theme primaryColorForView:self];
        } else {
            // normal and any other value
            c = UIButtonConfiguration.filledButtonConfiguration;
            c.baseBackgroundColor = [self.theme primaryColorForView:self];
        }

        // custom font (font, and size for info-small)
        c.titleTextAttributesTransformer = ^NSDictionary<NSAttributedStringKey, id> *_Nonnull(
            NSDictionary<NSAttributedStringKey, id> *_Nonnull incoming) {
            NSMutableDictionary<NSAttributedStringKey, id> *outgoing = [incoming mutableCopy];
            CGFloat fontSize = [DatamodelButtonStyleEnumInfoSmall isEqualToString:model.style]
                                   ? CM_SMALL_BUTTON_FONT_SIZE
                                   : CM_OS_BUTTON_FONT_SIZE;
            UIFont *font = [self.theme fontOfSize:fontSize];
            outgoing[NSFontAttributeName] = font;
            return outgoing;
        };

        return [UIButton buttonWithConfiguration:c primaryAction:nil];
    }

    return nil;
}

- (UIButton *)buildOldStyleButtonWithDataModel:(DatamodelButton *)model {
    // iOS 14 and earlier -- emulate iOS 15+ styles
    UIButton *button = [UIButton buttonWithType:UIButtonTypeCustom];
    button.layer.cornerRadius = 6.0;
    button.layer.masksToBounds = YES;
    button.contentEdgeInsets = UIEdgeInsetsMake(7, 0, 7, 0);

    UIColor *tintColor = [self.theme primaryColorForView:self];
    UIColor *backgroundColor = tintColor;

    CGFloat fontSize = CM_OS_BUTTON_FONT_SIZE;
    if ([DatamodelButtonStyleEnumLarge isEqualToString:model.style]) {
        button.contentEdgeInsets = UIEdgeInsetsMake(14, 0, 14, 0);
    } else if ([DatamodelButtonStyleEnumSecondary isEqualToString:model.style]) {
        // emulate iOS 15 tinted
        CGFloat h, s, b, a;
        [button setTitleColor:tintColor forState:UIControlStateNormal];
        [tintColor getHue:&h saturation:&s brightness:&b alpha:&a];
        // dark mode adjust tint for ios 13/14
        if (button.traitCollection.userInterfaceStyle == UIUserInterfaceStyleDark) {
            backgroundColor = [UIColor colorWithHue:h
                                         saturation:MAX(MIN(s * 0.84375, 1.0), 0.0)
                                         brightness:MAX(MIN(b * 0.3, 1.0), 0.0)
                                              alpha:1];
        } else {
            backgroundColor = [UIColor colorWithHue:h
                                         saturation:MAX(MIN(s * 0.10, 1.0), 0.0)
                                         brightness:MAX(MIN(b * 1.18, 1.0), 0.0)
                                              alpha:1];
        }
    } else if ([DatamodelButtonStyleEnumTertiary isEqualToString:model.style]) {
        [button setTitleColor:tintColor forState:UIControlStateNormal];
        if (@available(iOS 13.0, *)) {
            backgroundColor = [UIColor systemGray5Color];
        } else {
            backgroundColor = [UIColor colorWithRed:0.91 green:0.91 blue:0.91 alpha:1.0];
        }
    } else if ([DatamodelButtonStyleEnumInfo isEqualToString:model.style]) {
        backgroundColor = [UIColor clearColor];
        [button setTitleColor:tintColor forState:UIControlStateNormal];
    } else if ([DatamodelButtonStyleEnumInfoSmall isEqualToString:model.style]) {
        backgroundColor = [UIColor clearColor];
        [button setTitleColor:tintColor forState:UIControlStateNormal];
        fontSize = CM_SMALL_BUTTON_FONT_SIZE;
    } else {
        // normal and any other value
    }

    button.backgroundColor = backgroundColor;
    button.titleLabel.font = [self.theme fontOfSize:fontSize];

    return button;
}

@end
