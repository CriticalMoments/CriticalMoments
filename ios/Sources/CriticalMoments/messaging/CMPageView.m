//
//  CMPageView.m
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import "CMPageView.h"

#import "../themes/CMTheme_private.h"
#import "../utils/CMUtils.h"
#import "CMButton.h"
#import "CMGradientView.h"
#import "CMImageView.h"

#define CM_PAGE_SIDE_PADDING 25
#define CM_MIN_BTN_WIDTH 280
#define CM_SCROLL_SHIM_SIZE 20

@interface CMPageStack : NSObject

@property(nonatomic, readwrite) NSArray<UIView *> *views;
@property(nonatomic, readwrite) NSArray<NSNumber *> *spaceMultiplier;

@end
@implementation CMPageStack
@end

@interface CMPageView ()

@property(nonatomic, readwrite) CMTheme *customTheme;

@end

@implementation CMPageView

- (instancetype)initWithDatamodel:(DatamodelPage *)model andTheme:(CMTheme *)theme {
    self = [super init];
    if (self) {
        _customTheme = theme;
        [self buildSubviewsFromModel:model];
    }
    return self;
}

- (CMTheme *)theme {
    if (self.customTheme) {
        return self.customTheme;
    }
    return CMTheme.current;
}

- (void)buildSubviewsFromModel:(DatamodelPage *)model {
    self.backgroundColor = self.theme.backgroundColor;

    UIScrollView *scrollView = [[UIScrollView alloc] init];
    // add padding at the bottom
    scrollView.contentInset = UIEdgeInsetsMake(0, 0, CM_SCROLL_SHIM_SIZE, 0);
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

    CMGradientView *shimView = [[CMGradientView alloc] init];
    shimView.customTheme = self.theme;
    shimView.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:shimView];

    CMPageStack *stack = [self buildPageStack:model];
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
        [scrollView.bottomAnchor constraintEqualToAnchor:buttonArea.topAnchor],

        [topSpace.topAnchor constraintEqualToAnchor:scrollView.topAnchor],
        // important: frame, not the scroll view content or it might move
        [topSpace.heightAnchor constraintEqualToAnchor:self.heightAnchor multiplier:0.08],

        [stack.views.lastObject.bottomAnchor constraintEqualToAnchor:scrollView.bottomAnchor],

        [shimView.bottomAnchor constraintEqualToAnchor:buttonArea.topAnchor],
        [shimView.leftAnchor constraintEqualToAnchor:self.leftAnchor],
        [shimView.rightAnchor constraintEqualToAnchor:self.rightAnchor],
        [shimView.heightAnchor constraintEqualToConstant:CM_SCROLL_SHIM_SIZE],
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

        [constraints addObject:[view.topAnchor constraintEqualToSystemSpacingBelowAnchor:lastTop
                                                                              multiplier:topSpaceMultiplier]];

        if (view.frame.size.width == 0) {
            // Center with margin if width not explicitly set
            [constraints addObjectsFromArray:@[
                [view.topAnchor constraintEqualToSystemSpacingBelowAnchor:lastTop multiplier:topSpaceMultiplier],

                [view.leadingAnchor constraintEqualToAnchor:scrollView.layoutMarginsGuide.leadingAnchor
                                                   constant:CM_PAGE_SIDE_PADDING],
                [view.trailingAnchor constraintEqualToAnchor:scrollView.layoutMarginsGuide.trailingAnchor
                                                    constant:-CM_PAGE_SIDE_PADDING],
            ]];
        } else {
            // Center with exact width
            [constraints addObjectsFromArray:@[
                [view.widthAnchor constraintEqualToConstant:view.frame.size.width],
                [view.centerXAnchor constraintEqualToAnchor:scrollView.layoutMarginsGuide.centerXAnchor],
            ]];
        }

        lastTop = view.bottomAnchor;
    };

    NSArray<CMButton *> *buttons = [self buttons:model];
    if (buttons.count == 0) {
        [constraints addObjectsFromArray:@[
            [buttonArea.heightAnchor constraintEqualToConstant:0],
        ]];
    } else {
        NSLayoutYAxisAnchor *lastTop = buttonArea.topAnchor;
        for (UIButton *btn in buttons) {
            btn.translatesAutoresizingMaskIntoConstraints = NO;
            [buttonArea addSubview:btn];
            [constraints addObjectsFromArray:@[
                [btn.centerXAnchor constraintEqualToAnchor:buttonArea.centerXAnchor],
                [btn.topAnchor constraintEqualToSystemSpacingBelowAnchor:lastTop multiplier:1.6],

                // min width
                [btn.widthAnchor constraintGreaterThanOrEqualToConstant:CM_MIN_BTN_WIDTH],

                // max width
                [btn.leftAnchor constraintGreaterThanOrEqualToAnchor:buttonArea.leftAnchor],
                [btn.rightAnchor constraintLessThanOrEqualToAnchor:buttonArea.rightAnchor],
            ]];
            lastTop = btn.bottomAnchor;
        }

        // button above guide, and area
        CGFloat belowButtonPadding = buttons.count == 1 || [CMUtils isiPad] ? 30.0 : 0.0;
        [constraints addObjectsFromArray:@[
            // Last button up into layout margin guide. If only 1 button, move it up to more tapable position since we
            // have space
            [buttons.lastObject.bottomAnchor
                constraintGreaterThanOrEqualToAnchor:buttonArea.layoutMarginsGuide.bottomAnchor
                                            constant:-belowButtonPadding],

            [buttonArea.bottomAnchor constraintGreaterThanOrEqualToAnchor:buttons.lastObject.bottomAnchor],
        ]];
    }

    [NSLayoutConstraint activateConstraints:constraints];
}

- (CMPageStack *)buildPageStack:(DatamodelPage *)model {
    NSMutableArray<UIView *> *views = [[NSMutableArray alloc] init];
    NSMutableArray<NSNumber *> *spaces = [[NSMutableArray alloc] init];

    for (int i = 0; i < model.sectionCount; i++) {
        DatamodelPageSection *section = [model sectionAtIndex:i];
        if (!section)
            continue;

        UIView *sectionView = [self viewForSection:section];
        if (!sectionView)
            continue;

        [views addObject:sectionView];
        [spaces addObject:@(section.topSpacingScale)];
    }

    CMPageStack *stack = [[CMPageStack alloc] init];
    stack.views = views;
    stack.spaceMultiplier = spaces;
    return stack;
}

- (UIView *)viewForSection:(DatamodelPageSection *)section {
    if ([DatamodelSectionTypeEnumTitle isEqualToString:section.pageSectionType]) {
        // Titles
        return [self buildTitleView:section.titleData];
    } else if ([DatamodelSectionTypeEnumBodyText isEqualToString:section.pageSectionType]) {
        // Body
        return [self buildBodyView:section.bodyData];
    } else if ([DatamodelSectionTypeEnumImage isEqualToString:section.pageSectionType]) {
        // Image
        return [self buildImageView:section.imageData];
    }

    return nil;
}

- (UIView *)buildTitleView:(DatamodelTitlePageSection *)titleData {
    UILabel *titleView = [[UILabel alloc] init];

    titleView.text = titleData.title;
    titleView.numberOfLines = 0; // no limit

    if (titleData.centerText) {
        titleView.textAlignment = NSTextAlignmentCenter;
    }

    if (titleData.usePrimaryTextColor) {
        titleView.textColor = self.theme.primaryTextColor;
    } else {
        titleView.textColor = self.theme.secondaryTextColor;
    }

    CGFloat fontSize = self.theme.titleFontSize * titleData.scaleFactor;
    if (titleData.bold) {
        titleView.font = [self.theme boldFontOfSize:fontSize];
    } else {
        titleView.font = [self.theme fontOfSize:fontSize];
    }

    if (titleData.width != 0) {
        titleView.frame = CGRectMake(0, 0, titleData.width, 0);
    }

    return titleView;
}

- (UIView *)buildBodyView:(DatamodelBodyPageSection *)bodyData {
    UILabel *bodyLabel = [[UILabel alloc] init];
    bodyLabel.text = bodyData.bodyText;
    bodyLabel.numberOfLines = 0; // no limit

    if (bodyData.centerText) {
        bodyLabel.textAlignment = NSTextAlignmentCenter;
    }

    if (bodyData.usePrimaryTextColor) {
        bodyLabel.textColor = self.theme.primaryTextColor;
    } else {
        bodyLabel.textColor = self.theme.secondaryTextColor;
    }

    CGFloat fontSize = self.theme.bodyFontSize * bodyData.scaleFactor;
    if (bodyData.bold) {
        bodyLabel.font = [self.theme boldFontOfSize:fontSize];
    } else {
        bodyLabel.font = [self.theme fontOfSize:fontSize];
    }

    if (bodyData.width != 0) {
        bodyLabel.frame = CGRectMake(0, 0, bodyData.width, 0);
    }

    return bodyLabel;
}

- (UIView *)buildImageView:(DatamodelImage *)imageModel {
    CMImageView *iv = [[CMImageView alloc] initWithDatamodel:imageModel andTheme:self.theme];
    return iv;
}

- (NSArray<CMButton *> *)buttons:(DatamodelPage *)model {
    NSMutableArray<CMButton *> *buttons = [[NSMutableArray alloc] init];

    for (int i = 0; i < model.buttonsCount; i++) {
        DatamodelButton *buttonModel = [model buttonAtIndex:i];
        if (!buttonModel)
            continue;

        CMButton *button = [[CMButton alloc] initWithWithDataModel:buttonModel andTheme:self.theme];
        if (button) {
            [buttons addObject:button];
            __weak CMPageView *weakSelf = self;
            button.defaultAction = ^{
              if (weakSelf.anyButtonDefaultAction) {
                  weakSelf.anyButtonDefaultAction();
              }
            };
            button.buttonTappedAction = ^{
              if (weakSelf.buttonCallback) {
                  weakSelf.buttonCallback(buttonModel.title, i);
              }
            };
        }
    }

    return buttons;
}

@end
