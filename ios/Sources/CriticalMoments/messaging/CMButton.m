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

@implementation CMButton

+ (UIButton *)buttonWithWithDataModel:(NSObject *)o andTheme:(CMTheme *)theme {
    UIButton *button;

    if (@available(iOS 15.0, *)) {
        // TODO: real data model prop
        NSString *buttonType = o;
        UIButtonConfiguration *c;
        if ([@"large" isEqualToString:buttonType]) {
            c = UIButtonConfiguration.filledButtonConfiguration;
            c.buttonSize = UIButtonConfigurationSizeLarge;
        } else if ([@"secondary" isEqualToString:buttonType]) {
            c = UIButtonConfiguration.tintedButtonConfiguration;
        } else if ([@"tertiary" isEqualToString:buttonType]) {
            c = UIButtonConfiguration.grayButtonConfiguration;
        } else if ([@"info" isEqualToString:buttonType]) {
            c = UIButtonConfiguration.plainButtonConfiguration;
        } else if ([@"info-small" isEqualToString:buttonType]) {
            c = UIButtonConfiguration.plainButtonConfiguration;
        } else {
            // normal and any other value
            c = UIButtonConfiguration.filledButtonConfiguration;
        }

        // custom font (font, and size for info-small)
        NSLog(@"");
        c.titleTextAttributesTransformer = ^NSDictionary<NSAttributedStringKey, id> *_Nonnull(
            NSDictionary<NSAttributedStringKey, id> *_Nonnull incoming) {
            NSMutableDictionary<NSAttributedStringKey, id> *outgoing = [incoming mutableCopy];
            CGFloat fontSize =
                [@"info-small" isEqualToString:buttonType] ? CM_SMALL_BUTTON_FONT_SIZE : CM_OS_BUTTON_FONT_SIZE;
            outgoing[NSFontAttributeName] = [theme fontOfSize:fontSize];
            return outgoing;
        };

        button = [UIButton buttonWithConfiguration:c primaryAction:nil];

        // TODO: font? System really fighting me on changing button fonts.
        /*theme = [[CMTheme alloc] init];
        theme.fontName = @"Avenir-Black";
        button.titleLabel.font = [theme fontOfSize:UIFont.systemFontSize];
        [button.titleLabel setFont:[UIFont fontWithName:@"Helvetica-Bold" size:39.0]];*/
    } else {
        // iOS 14 and earlier -- emulate iOS 15+ styles
        button = [UIButton buttonWithType:UIButtonTypeCustom];
        button.layer.cornerRadius = 6.0;
        button.layer.masksToBounds = YES;
        button.contentEdgeInsets = UIEdgeInsetsMake(7, 0, 7, 0);

        UIColor *tintColor = [theme primaryColorForView:button];
        UIColor *backgroundColor = tintColor;

        // TODO: real data model prop
        NSString *buttonType = o;
        CGFloat fontSize = CM_OS_BUTTON_FONT_SIZE;
        if ([@"large" isEqualToString:buttonType]) {
            button.contentEdgeInsets = UIEdgeInsetsMake(14, 0, 14, 0);
        } else if ([@"secondary" isEqualToString:buttonType]) {
            // emulate iOS 15 tinted
            CGFloat h, s, b, a;
            [button setTitleColor:tintColor forState:UIControlStateNormal];
            [tintColor getHue:&h saturation:&s brightness:&b alpha:&a];
            // dark mode adjust tint for ios 13/14
            if (button.traitCollection.userInterfaceStyle == UIUserInterfaceStyleDark) {
                backgroundColor = [UIColor colorWithHue:h
                                             saturation:MAX(MIN(s - 0.12, 1.0), 0.0)
                                             brightness:MAX(MIN(b - 0.6, 1.0), 0.0)
                                                  alpha:1];
            } else {
                backgroundColor = [UIColor colorWithHue:h
                                             saturation:MAX(s - 0.52, 0.0)
                                             brightness:MAX(MIN(b + 0.1, 1.0), 0.0)
                                                  alpha:1];
            }
        } else if ([@"tertiary" isEqualToString:buttonType]) {
            [button setTitleColor:tintColor forState:UIControlStateNormal];
            if (@available(iOS 13.0, *)) {
                backgroundColor = [UIColor systemGray5Color];
            } else {
                backgroundColor = [UIColor colorWithRed:0.91 green:0.91 blue:0.91 alpha:1.0];
            }
        } else if ([@"info" isEqualToString:buttonType]) {
            backgroundColor = [UIColor clearColor];
            [button setTitleColor:tintColor forState:UIControlStateNormal];
        } else if ([@"info-small" isEqualToString:buttonType]) {
            backgroundColor = [UIColor clearColor];
            [button setTitleColor:tintColor forState:UIControlStateNormal];
            fontSize = CM_SMALL_BUTTON_FONT_SIZE;
        } else {
            // normal and any other value
        }

        button.backgroundColor = backgroundColor;
        button.titleLabel.font = [theme fontOfSize:fontSize];
    }

    // TODO: real value
    [button setTitle:o forState:UIControlStateNormal];

    return button;
}

@end
