//
//  InfoHeader.m
//  SampleApp
//
//  Created by Steve Cosman on 2024-02-02.
//

#import "InfoHeader.h"

#import "Utils.h"

@implementation InfoHeader

+ (instancetype)headerWithScreen:(CMDemoScreen *)screen {
    InfoHeader *header = [[InfoHeader alloc] init];
    header.preservesSuperviewLayoutMargins = YES;
    header.autoresizingMask = UIViewAutoresizingFlexibleHeight;
    NSMutableArray<NSLayoutConstraint *> *constraints = [[NSMutableArray alloc] init];
    UIView *lastItem;

    if (screen.infoText) {
        UILabel *infoLabel = [[UILabel alloc] init];
        infoLabel.translatesAutoresizingMaskIntoConstraints = NO;
        [header addSubview:infoLabel];
        infoLabel.numberOfLines = 0;
        infoLabel.text = screen.infoText;
        if (@available(iOS 13.0, *)) {
            infoLabel.textColor = [UIColor secondaryLabelColor];
        } else {
            infoLabel.textColor = [UIColor systemGrayColor];
        }
        infoLabel.font = [UIFont boldSystemFontOfSize:UIFont.systemFontSize];

        [constraints addObjectsFromArray:@[
            [infoLabel.topAnchor constraintEqualToAnchor:header.layoutMarginsGuide.topAnchor],
            [infoLabel.leadingAnchor constraintEqualToAnchor:header.layoutMarginsGuide.leadingAnchor],
            [infoLabel.trailingAnchor constraintEqualToAnchor:header.layoutMarginsGuide.trailingAnchor],
        ]];

        lastItem = infoLabel;
    }

    if (@available(iOS 14.0, *)) {
        if (screen.buttonLink) {
            UIButton *button = [UIButton buttonWithType:UIButtonTypeSystem];
            [button setTitle:(screen.buttonTitle ? screen.buttonTitle : @"See the docs") forState:UIControlStateNormal];
            button.translatesAutoresizingMaskIntoConstraints = NO;
            __block NSString *link = screen.buttonLink;
            [header addSubview:button];
            UIAction *linkAction = [UIAction actionWithHandler:^(UIAction *action) {
              [[UIApplication sharedApplication] openURL:[NSURL URLWithString:link] options:@{} completionHandler:nil];
            }];
            [button addAction:linkAction forControlEvents:UIControlEventTouchUpInside];

            if (lastItem) {
                [constraints addObject:[button.topAnchor constraintEqualToSystemSpacingBelowAnchor:lastItem.bottomAnchor
                                                                                        multiplier:1]];
            } else {
                [constraints
                    addObject:[button.topAnchor constraintEqualToAnchor:header.layoutMarginsGuide.bottomAnchor]];
            }

            [constraints addObjectsFromArray:@[
                [button.leadingAnchor constraintEqualToAnchor:header.layoutMarginsGuide.leadingAnchor],
            ]];
            lastItem = button;
        }
    }

    // Color bars
    NSArray<UIColor *> *barColors = @[
        [Utils colorFromHexString:@"#EBA239"], [Utils colorFromHexString:@"#DF612F"],
        [Utils colorFromHexString:@"#D73A28"], [Utils colorFromHexString:@"#9F2F42"],
        [Utils colorFromHexString:@"#60233E"]
    ];
    UIView *lastBar;
    for (UIColor *barColor in barColors) {
        UIView *barView = [[UIView alloc] init];
        barView.backgroundColor = barColor;
        barView.translatesAutoresizingMaskIntoConstraints = NO;
        [header addSubview:barView];

        NSLayoutConstraint *topAnchor;
        if (lastBar) {
            topAnchor = [barView.topAnchor constraintEqualToAnchor:lastBar.bottomAnchor];
        } else if (lastItem) {
            topAnchor = [barView.topAnchor constraintEqualToSystemSpacingBelowAnchor:lastItem.bottomAnchor
                                                                          multiplier:1.5];
        } else {
            topAnchor = [barView.topAnchor constraintEqualToAnchor:header.topAnchor];
        }

        [constraints addObjectsFromArray:@[
            topAnchor, [barView.leadingAnchor constraintEqualToAnchor:header.leadingAnchor],
            [barView.trailingAnchor constraintEqualToAnchor:header.trailingAnchor],
            [barView.heightAnchor constraintEqualToConstant:7.0]
        ]];

        lastBar = barView;
    }
    lastItem = lastBar;

    [constraints addObject:[lastItem.bottomAnchor constraintEqualToAnchor:header.layoutMarginsGuide.bottomAnchor]];
    [NSLayoutConstraint activateConstraints:constraints];

    return header;
}

@end
