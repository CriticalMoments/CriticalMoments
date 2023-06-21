//
//  CMSheetViewController.m
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import "CMModalViewController.h"

#import "../themes/CMTheme_private.h"
#import "CMPageView.h"

@interface CMModalViewController ()

@property(nonnull, strong, readwrite) DatamodelModalAction *model;
@property(nonnull, strong, readwrite) UIButton *closeButton;
@property(nonatomic, readwrite) CMTheme *customTheme;

@end

@implementation CMModalViewController

- (instancetype)initWithDatamodel:(DatamodelModalAction *)model {
    self = [super init];
    if (self) {
        self.model = model;
        if (model.customThemeName.length > 0) {
            CMTheme *customTheme = [CMTheme namedThemeFromAppcore:model.customThemeName];
            self.customTheme = customTheme;
        }

        // prevents swipe to dismiss
        if (@available(iOS 13.0, *)) {
            self.modalInPresentation = !model.showCloseButton;
        } else {
            // TODO P0 confirm: sheets are new in 13 so might be no-op
        }
    }
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];

    // Theme
    CMTheme *theme = self.customTheme;
    if (!theme) {
        theme = CMTheme.current;
    }
    self.view.backgroundColor = theme.backgroundColor;

    CMPageView *pv = [[CMPageView alloc] initWithDatamodel:self.model.content andTheme:theme];
    __weak CMModalViewController *weakSelf = self;
    pv.anyButtonDefaultAction = ^{
      [weakSelf dismissSheet];
    };
    pv.translatesAutoresizingMaskIntoConstraints = NO;
    [self.view addSubview:pv];

    UIButton *closeBtn;
    if (@available(iOS 13.0, *)) {
        closeBtn = [UIButton buttonWithType:UIButtonTypeClose];
    } else {
        closeBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        // TODO: check this unicode on earliest deployment target: ios 12
        [closeBtn setTitle:@"✕" forState:UIControlStateNormal];
        [closeBtn setTitleColor:theme.primaryTextColor forState:UIControlStateNormal];
    }
    closeBtn.hidden = !self.model.showCloseButton;
    closeBtn.translatesAutoresizingMaskIntoConstraints = NO;
    [closeBtn addTarget:self
                  action:@selector(closeButtonTapped:)
        forControlEvents:UIControlEventPrimaryActionTriggered];
    [self.view addSubview:closeBtn];
    _closeButton = closeBtn;

    // Layout

    NSArray<NSLayoutConstraint *> *constraints = @[
        [closeBtn.topAnchor constraintEqualToSystemSpacingBelowAnchor:self.view.topAnchor multiplier:2.0],
        [closeBtn.rightAnchor constraintEqualToAnchor:self.view.layoutMarginsGuide.rightAnchor],

        [pv.topAnchor constraintEqualToAnchor:self.view.topAnchor],
        [pv.leftAnchor constraintEqualToAnchor:self.view.leftAnchor],
        [pv.rightAnchor constraintEqualToAnchor:self.view.rightAnchor],
        [pv.bottomAnchor constraintEqualToAnchor:self.view.bottomAnchor],
    ];
    [NSLayoutConstraint activateConstraints:constraints];
}

- (void)closeButtonTapped:(UIButton *)sender {
    [self dismissSheet];
}

- (void)dismissSheet {
    dispatch_async(dispatch_get_main_queue(), ^{
      UIViewController *pvc = self.presentingViewController;
      [pvc dismissViewControllerAnimated:YES completion:nil];
    });
}

@end
