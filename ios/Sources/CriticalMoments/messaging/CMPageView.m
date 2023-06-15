//
//  CMPageView.m
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import "CMPageView.h"

#import "CMImageView.h"

#define CM_PAGE_SIDE_PADDING 40
#define CM_SIMPLE_IMAGE_SIZE 40

@implementation CMPageView

- (instancetype)init {
    self = [super init];
    if (self) {
        [self buildSubviews];
    }
    return self;
}

- (void)buildSubviews {
    // TODO case for layouts
    [self buildSubviewsForSimpleLayout];
}

- (void)buildSubviewsForSimpleLayout {
    // Theme
    CMTheme *theme = self.customTheme;
    if (!theme) {
        theme = CMTheme.current;
    }

    self.backgroundColor = theme.backgroundColor;

    UIScrollView *scrollView = [[UIScrollView alloc] init];
    // add padding at the bottom
    scrollView.contentInset = UIEdgeInsetsMake(0, 0, 12, 0);
    scrollView.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:scrollView];

    // UIScrollView doesn't like UILayoutGuide
    UIView *topSpace = [[UIView alloc] init];
    topSpace.translatesAutoresizingMaskIntoConstraints = NO;
    [scrollView addSubview:topSpace];

    CMImageView *iv = [[CMImageView alloc] init];
    iv.translatesAutoresizingMaskIntoConstraints = NO;
    [scrollView addSubview:iv];

    UILabel *titleView = [[UILabel alloc] init];
    // TODO
    titleView.text = @"Important Announcement!";
    titleView.numberOfLines = 0; // no limit
    titleView.textAlignment = NSTextAlignmentCenter;
    titleView.textColor = theme.primaryTextColor;
    titleView.font = [theme boldFontOfSize:theme.titleFontSize];
    titleView.translatesAutoresizingMaskIntoConstraints = NO;
    [scrollView addSubview:titleView];

    UILabel *subtitle = [[UILabel alloc] init];
    // TODO
    subtitle.text = @"New pricing coming soon.";
    subtitle.numberOfLines = 0; // no limit
    subtitle.textAlignment = NSTextAlignmentCenter;
    subtitle.textColor = theme.primaryTextColor;
    subtitle.font = [theme boldFontOfSize:theme.subtitleFontSize];
    subtitle.translatesAutoresizingMaskIntoConstraints = NO;
    [scrollView addSubview:subtitle];

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
    bodyLabel.textColor = theme.secondaryTextColor;
    bodyLabel.font = [theme fontOfSize:UIFont.systemFontSize];
    bodyLabel.translatesAutoresizingMaskIntoConstraints = NO;
    [scrollView addSubview:bodyLabel];

    // Layout

    NSArray<NSLayoutConstraint *> *constraints = @[
        [scrollView.topAnchor constraintEqualToAnchor:self.topAnchor],
        [scrollView.leftAnchor constraintEqualToAnchor:self.leftAnchor],
        // TODO
        [scrollView.bottomAnchor constraintEqualToAnchor:self.bottomAnchor],
        [scrollView.rightAnchor constraintEqualToAnchor:self.rightAnchor],

        [topSpace.topAnchor constraintEqualToAnchor:scrollView.topAnchor],
        // important: frame, not the scroll view content or it might move
        [topSpace.heightAnchor constraintEqualToAnchor:self.heightAnchor multiplier:0.08],

        [iv.topAnchor constraintEqualToAnchor:topSpace.bottomAnchor],
        [iv.centerXAnchor constraintEqualToAnchor:scrollView.centerXAnchor],
        [iv.heightAnchor constraintEqualToConstant:CM_SIMPLE_IMAGE_SIZE],
        [iv.widthAnchor constraintEqualToConstant:CM_SIMPLE_IMAGE_SIZE],

        [titleView.topAnchor constraintEqualToSystemSpacingBelowAnchor:iv.bottomAnchor multiplier:2.0],
        //[titleView.topAnchor constraintEqualToAnchor:topSpace.bottomAnchor],
        [titleView.leadingAnchor constraintGreaterThanOrEqualToAnchor:scrollView.layoutMarginsGuide.leadingAnchor
                                                             constant:CM_PAGE_SIDE_PADDING],
        [titleView.trailingAnchor constraintLessThanOrEqualToAnchor:scrollView.layoutMarginsGuide.trailingAnchor
                                                           constant:-CM_PAGE_SIDE_PADDING],
        [titleView.widthAnchor constraintLessThanOrEqualToConstant:CM_MAX_TEXT_WIDTH],
        [titleView.centerXAnchor constraintEqualToAnchor:scrollView.centerXAnchor],

        [subtitle.topAnchor constraintEqualToSystemSpacingBelowAnchor:titleView.bottomAnchor multiplier:2.0],
        //[titleView.topAnchor constraintEqualToAnchor:topSpace.bottomAnchor],
        [subtitle.leadingAnchor constraintGreaterThanOrEqualToAnchor:scrollView.layoutMarginsGuide.leadingAnchor
                                                            constant:CM_PAGE_SIDE_PADDING],
        [subtitle.trailingAnchor constraintLessThanOrEqualToAnchor:scrollView.layoutMarginsGuide.trailingAnchor
                                                          constant:-CM_PAGE_SIDE_PADDING],
        [subtitle.widthAnchor constraintLessThanOrEqualToConstant:CM_MAX_TEXT_WIDTH],
        [subtitle.centerXAnchor constraintEqualToAnchor:scrollView.centerXAnchor],

        [bodyLabel.topAnchor constraintEqualToSystemSpacingBelowAnchor:subtitle.bottomAnchor multiplier:4.0],
        [bodyLabel.leadingAnchor constraintGreaterThanOrEqualToAnchor:scrollView.layoutMarginsGuide.leadingAnchor
                                                             constant:CM_PAGE_SIDE_PADDING],
        [bodyLabel.trailingAnchor constraintLessThanOrEqualToAnchor:scrollView.layoutMarginsGuide.trailingAnchor
                                                           constant:-CM_PAGE_SIDE_PADDING],
        [bodyLabel.widthAnchor constraintLessThanOrEqualToConstant:CM_MAX_TEXT_WIDTH],
        [bodyLabel.centerXAnchor constraintEqualToAnchor:scrollView.centerXAnchor],

        [bodyLabel.bottomAnchor constraintEqualToAnchor:scrollView.bottomAnchor],
    ];
    [NSLayoutConstraint activateConstraints:constraints];
}

@end
