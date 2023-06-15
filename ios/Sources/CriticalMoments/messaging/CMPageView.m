//
//  CMPageView.m
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import "CMPageView.h"

#import "CMButton.h"
#import "CMImageView.h"

#define CM_PAGE_SIDE_PADDING 40
#define CM_SIMPLE_IMAGE_SIZE 40
#define CM_MIN_BTN_WIDTH 240

@interface CMPageStack : NSObject

@property(nonatomic, readwrite) NSArray<UIView *> *views;
@property(nonatomic, readwrite) NSArray<NSNumber *> *spaceMultiplier;

@end
@implementation CMPageStack
@end

@implementation CMPageView

- (instancetype)init {
    self = [super init];
    if (self) {
        [self buildSubviews];
    }
    return self;
}

- (CMTheme *)theme {
    if (self.customTheme) {
        return self.customTheme;
    }
    return CMTheme.current;
}

- (void)buildSubviews {
    // TODO case for layouts
    [self buildSubviewsForSimpleLayout];
}

- (void)buildSubviewsForSimpleLayout {
    self.backgroundColor = self.theme.backgroundColor;

    UIScrollView *scrollView = [[UIScrollView alloc] init];
    // add padding at the bottom
    scrollView.contentInset = UIEdgeInsetsMake(0, 0, 12, 0);
    scrollView.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:scrollView];

    // UIScrollView doesn't like UILayoutGuide
    UIView *topSpace = [[UIView alloc] init];
    topSpace.translatesAutoresizingMaskIntoConstraints = NO;
    [scrollView addSubview:topSpace];

    UIView *buttonArea = [[UIView alloc] init];
    buttonArea.backgroundColor = self.theme.backgroundColor;
    buttonArea.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:buttonArea];

    CMPageStack *stack = self.simpleLayoutViewStack;
    if (stack.views.count == 0) {
        return;
    }

    // Layout

    NSMutableArray<NSLayoutConstraint *> *constraints = [[NSMutableArray alloc] initWithArray:@[
        [buttonArea.bottomAnchor constraintEqualToAnchor:self.bottomAnchor],
        [buttonArea.leftAnchor constraintEqualToAnchor:self.leftAnchor],
        [buttonArea.rightAnchor constraintEqualToAnchor:self.rightAnchor],

        [scrollView.topAnchor constraintEqualToAnchor:self.topAnchor],
        [scrollView.leftAnchor constraintEqualToAnchor:self.leftAnchor],
        [scrollView.rightAnchor constraintEqualToAnchor:self.rightAnchor],
        // TODO
        [scrollView.bottomAnchor constraintEqualToAnchor:buttonArea.topAnchor],

        [topSpace.topAnchor constraintEqualToAnchor:scrollView.topAnchor],
        // important: frame, not the scroll view content or it might move
        [topSpace.heightAnchor constraintEqualToAnchor:self.heightAnchor multiplier:0.08],

        [stack.views.lastObject.bottomAnchor constraintEqualToAnchor:scrollView.bottomAnchor],
    ]];

    NSLayoutYAxisAnchor *lastTop = topSpace.bottomAnchor;
    for (int i = 0; i < stack.views.count; i++) {
        UIView *view = [stack.views objectAtIndex:i];
        view.translatesAutoresizingMaskIntoConstraints = NO;
        [scrollView addSubview:view];

        CGFloat topSpaceMultiplier = 0.0;
        if (i < stack.spaceMultiplier.count) {
            topSpaceMultiplier = [stack.spaceMultiplier objectAtIndex:i].floatValue;
        }

        [constraints addObjectsFromArray:@[
            [view.topAnchor constraintEqualToSystemSpacingBelowAnchor:lastTop multiplier:topSpaceMultiplier],
            [view.centerXAnchor constraintEqualToAnchor:scrollView.centerXAnchor],
            [view.leadingAnchor constraintGreaterThanOrEqualToAnchor:scrollView.layoutMarginsGuide.leadingAnchor
                                                            constant:CM_PAGE_SIDE_PADDING],
            [view.trailingAnchor constraintLessThanOrEqualToAnchor:scrollView.layoutMarginsGuide.trailingAnchor
                                                          constant:-CM_PAGE_SIDE_PADDING],
            [view.widthAnchor constraintLessThanOrEqualToConstant:CM_MAX_TEXT_WIDTH],
        ]];

        lastTop = view.bottomAnchor;
    };

    NSArray<UIButton *> *buttons = [self buttons];
    if (buttons.count == 0) {
        // TODO test this
        [constraints addObjectsFromArray:@[
            [buttonArea.heightAnchor constraintEqualToConstant:0],
        ]];
    } else {
        // TODO layoutGuide
        NSLayoutYAxisAnchor *lastTop = buttonArea.topAnchor;
        for (UIButton *btn in buttons) {
            btn.translatesAutoresizingMaskIntoConstraints = NO;
            [buttonArea addSubview:btn];
            [constraints addObjectsFromArray:@[
                [btn.centerXAnchor constraintEqualToAnchor:buttonArea.centerXAnchor],
                [btn.topAnchor constraintEqualToSystemSpacingBelowAnchor:lastTop multiplier:2.0],

                // min width
                [btn.widthAnchor constraintGreaterThanOrEqualToConstant:CM_MIN_BTN_WIDTH],

                // max width
                [btn.leftAnchor constraintGreaterThanOrEqualToAnchor:buttonArea.leftAnchor],
                [btn.rightAnchor constraintLessThanOrEqualToAnchor:buttonArea.rightAnchor],
            ]];
            lastTop = btn.bottomAnchor;
        }

        // button above guide, and area
        [constraints addObjectsFromArray:@[
            [buttons.lastObject.bottomAnchor
                constraintGreaterThanOrEqualToAnchor:buttonArea.layoutMarginsGuide.bottomAnchor],
            [buttonArea.bottomAnchor constraintGreaterThanOrEqualToAnchor:buttons.lastObject.bottomAnchor],
        ]];
    }

    [NSLayoutConstraint activateConstraints:constraints];
}

- (CMPageStack *)simpleLayoutViewStack {
    CMImageView *iv = [[CMImageView alloc] init];
    NSArray<NSLayoutConstraint *> *constraints = @[
        [iv.heightAnchor constraintEqualToConstant:CM_SIMPLE_IMAGE_SIZE],
        [iv.widthAnchor constraintEqualToConstant:CM_SIMPLE_IMAGE_SIZE],
    ];
    [NSLayoutConstraint activateConstraints:constraints];

    UILabel *titleView = [[UILabel alloc] init];
    // TODO
    titleView.text = @"Important Announcement!";
    titleView.numberOfLines = 0; // no limit
    titleView.textAlignment = NSTextAlignmentCenter;
    titleView.textColor = self.theme.primaryTextColor;
    titleView.font = [self.theme boldFontOfSize:self.theme.titleFontSize];

    UILabel *subtitle = [[UILabel alloc] init];
    // TODO
    subtitle.text = @"New pricing coming soon.";
    subtitle.numberOfLines = 0; // no limit
    subtitle.textAlignment = NSTextAlignmentCenter;
    subtitle.textColor = self.theme.primaryTextColor;
    subtitle.font = [self.theme boldFontOfSize:self.theme.subtitleFontSize];

    UILabel *bodyLabel = [[UILabel alloc] init];
    // TODO
    bodyLabel.text =
        @"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris nec eros imperdiet, ullamcorper neque "
        @"sagittis, ultrices lacus. Praesent sed mattis odio, a feugiat risus. Aliquam erat volutpat. Quisque "
        @"condimentum leo sapien, nec ullamcorper diam viverra sollicitudin. Nam id laoreet metus. Ut quis maximus "
        @"lacus. Aliquam sodales dui quis leo ullamcorper porta. Sed sed varius enim. Vivamus viverra consectetur erat "
        @"sit amet eleifend. Morbi facilisis, dolor a placerat rhoncus, lorem dolor egestas mauris, vel tincidunt sem "
        @"nisi a ante. Donec eget elit rhoncus, ornare leo vel, pharetra ipsum. Maecenas at dolor non dolor feugiat "
        @"condimentum. Suspendisse nec dignissim magna. Maecenas at sem eu ante egestas tempor. Duis enim enim, "
        @"faucibus sed turpis sit amet, ultricies feugiat turpis. Quisque posuere sagittis mauris, et tempus diam "
        @"maximus vel. Curabitur eget quam eu nisl elementum tincidunt faucibus in eros. Morbi suscipit lorem nisi, at "
        @"dapibus mi elementum vitae. Suspendisse tempus maximus diam sed blandit. Fusce dignissim velit at dapibus "
        @"rutrum. Vestibulum sollicitudin nec nunc at molestie. Mauris vel nisi tincidunt, efficitur velit a, congue "
        @"odio. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris nec eros imperdiet, ullamcorper neque "
        @"sagittis, ultrices lacus. Praesent sed mattis odio, a feugiat risus. Aliquam erat volutpat. Quisque "
        @"condimentum leo sapien, nec ullamcorper diam viverra sollicitudin. Nam id laoreet metus. Ut quis maximus "
        @"lacus. Aliquam sodales dui quis leo ullamcorper porta. Sed sed varius enim. Vivamus viverra consectetur erat "
        @"sit amet eleifend. Morbi facilisis, dolor a placerat rhoncus, lorem dolor egestas mauris, vel tincidunt sem "
        @"nisi a ante. Donec eget elit rhoncus, ornare leo vel, pharetra ipsum. Maecenas at dolor non dolor feugiat "
        @"condimentum. Suspendisse nec dignissim magna. Maecenas at sem eu ante egestas tempor. Duis enim enim, "
        @"faucibus sed turpis sit amet, ultricies feugiat turpis. Quisque posuere sagittis mauris, et tempus diam "
        @"maximus vel. Curabitur eget quam eu nisl elementum tincidunt faucibus in eros. Morbi suscipit lorem nisi, at "
        @"dapibus mi elementum vitae. Suspendisse tempus maximus diam sed blandit. Fusce dignissim velit at dapibus "
        @"rutrum. Vestibulum sollicitudin nec nunc at molestie. Mauris vel nisi tincidunt, efficitur velit a, congue "
        @"odio.";
    bodyLabel.numberOfLines = 0; // no limit
    bodyLabel.textAlignment = NSTextAlignmentCenter;
    bodyLabel.textColor = self.theme.secondaryTextColor;
    bodyLabel.font = [self.theme fontOfSize:UIFont.systemFontSize];

    CMPageStack *stack = [[CMPageStack alloc] init];
    stack.views = @[ iv, titleView, subtitle, bodyLabel ];
    stack.spaceMultiplier = @[ @0, @1, @1.5, @4 ];
    return stack;
}

- (NSArray<UIButton *> *)buttons {
    UIButton *primary = [CMButton buttonWithWithDataModel:@"large" andTheme:self.theme];    // large fill
    UIButton *secondary = [CMButton buttonWithWithDataModel:@"normal" andTheme:self.theme]; // fill
    UIButton *third = [CMButton buttonWithWithDataModel:@"secondary" andTheme:self.theme];  // tinted
    UIButton *forth = [CMButton buttonWithWithDataModel:@"tertiary" andTheme:self.theme];   // grey
    UIButton *fifth = [CMButton buttonWithWithDataModel:@"info" andTheme:self.theme];       // text
    UIButton *sixth = [CMButton buttonWithWithDataModel:@"info-small" andTheme:self.theme]; // text

    return @[ primary, secondary, third, forth, fifth, sixth ];
}

@end
