//
//  CMPageView.m
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import "CMPageView.h"

#import "CMButton.h"
#import "CMGradientView.h"
#import "CMImageView.h"

#define CM_PAGE_SIDE_PADDING 40
#define CM_MIN_BTN_WIDTH 280
#define CM_SCROLL_SHIM_SIZE 20

@interface CMPageStack : NSObject

@property(nonatomic, readwrite) NSArray<UIView *> *views;
@property(nonatomic, readwrite) NSArray<NSNumber *> *spaceMultiplier;

@end
@implementation CMPageStack
@end

@implementation CMPageView

- (instancetype)initWithDatamodel:(DatamodelPage *)model {
    self = [super init];
    if (self) {
        [self buildSubviewsFromModel:model];
    }
    return self;
}

// TODO: if set after init, will this have any effect??
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
    shimView.customTheme = self.customTheme;
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
        // TODO
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

        [constraints addObjectsFromArray:@[
            [view.topAnchor constraintEqualToSystemSpacingBelowAnchor:lastTop multiplier:topSpaceMultiplier],
            // TODO: short non-centered text not working
            [view.centerXAnchor constraintEqualToAnchor:scrollView.centerXAnchor],
            [view.leadingAnchor constraintGreaterThanOrEqualToAnchor:scrollView.layoutMarginsGuide.leadingAnchor
                                                            constant:CM_PAGE_SIDE_PADDING],
            [view.trailingAnchor constraintLessThanOrEqualToAnchor:scrollView.layoutMarginsGuide.trailingAnchor
                                                          constant:-CM_PAGE_SIDE_PADDING],
            [view.widthAnchor constraintLessThanOrEqualToConstant:CM_MAX_TEXT_WIDTH],
        ]];

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
        CGFloat belowButtonPadding = buttons.count > 1 ? 0.0 : 30.0;
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
    // Titles
    if ([DatamodelSectionTypeEnumTitle isEqualToString:section.pageSectionType]) {
        return [self buildTitleView:section.titleData];
    }

    // Body
    if ([DatamodelSectionTypeEnumBodyText isEqualToString:section.pageSectionType]) {
        return [self buildBodyView:section.bodyData];
    }

    // Image
    if ([DatamodelSectionTypeEnumImage isEqualToString:section.pageSectionType]) {
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
    titleView.textColor = self.theme.primaryTextColor;

    CGFloat fontSize = self.theme.titleFontSize * titleData.scaleFactor;
    if (titleData.bold) {
        titleView.font = [self.theme boldFontOfSize:fontSize];
    } else {
        titleView.font = [self.theme fontOfSize:fontSize];
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

    CGFloat fontSize = UIFont.systemFontSize * bodyData.scaleFactor;
    if (bodyData.bold) {
        bodyLabel.font = [self.theme boldFontOfSize:fontSize];
    } else {
        bodyLabel.font = [self.theme fontOfSize:fontSize];
    }

    return bodyLabel;
}

- (UIView *)buildImageView:(DatamodelImage *)imageModel {
    CMImageView *iv = [[CMImageView alloc] initWithDatamodel:imageModel];
    return iv;
}

- (NSArray<CMButton *> *)buttons:(DatamodelPage *)model {
    NSMutableArray<CMButton *> *buttons = [[NSMutableArray alloc] init];

    // TODO: preventDefault
    // TODO: actionName

    for (int i = 0; i < model.buttonsCount; i++) {
        DatamodelButton *buttonModel = [model buttonAtIndex:i];
        if (!buttonModel)
            continue;

        CMButton *button = [[CMButton alloc] initWithWithDataModel:buttonModel andTheme:self.customTheme];
        if (button) {
            [buttons addObject:button];
            __weak CMPageView *weakSelf = self;
            button.defaultAction = ^{
              if (weakSelf.anyButtonDefaultAction) {
                  weakSelf.anyButtonDefaultAction();
              }
            };
        }
    }

    return buttons;
}

@end
