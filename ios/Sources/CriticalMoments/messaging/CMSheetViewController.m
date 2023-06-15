//
//  CMSheetViewController.m
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import "CMSheetViewController.h"

#import "CMPageView.h"

@interface CMSheetViewController ()

@property(nonnull, strong, readwrite) UIButton *closeButton;

@end

@implementation CMSheetViewController

- (instancetype)init {
    self = [super init];
    if (self) {
        self.allowDismissing = YES;
    }
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.

    // Theme
    CMTheme *theme = self.customTheme;
    if (!theme) {
        theme = CMTheme.current;
    }

    // TODO Theme
    self.view.backgroundColor = theme.backgroundColor;

    UIButton *closeBtn;
    if (@available(iOS 13.0, *)) {
        closeBtn = [UIButton buttonWithType:UIButtonTypeClose];
    } else {
        closeBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        // TODO: check this unicode on earliest deployment target: ios 12
        [closeBtn setTitle:@"✕" forState:UIControlStateNormal];
        [closeBtn setTitleColor:theme.primaryTextColor forState:UIControlStateNormal];
    }
    closeBtn.hidden = !_allowDismissing;
    closeBtn.translatesAutoresizingMaskIntoConstraints = NO;
    [closeBtn addTarget:self
                  action:@selector(closeButtonTapped:)
        forControlEvents:UIControlEventPrimaryActionTriggered];
    [self.view addSubview:closeBtn];
    _closeButton = closeBtn;

    // TODO grabber?

    // TODO tint color from code
    // https://developer.apple.com/documentation/uikit/uiview/1622467-tintcolor?language=objc

    // TODO
    // self.sheetPresentationController.detents = @[UISheetPresentationControllerDetentIdentifierLarge];

    CMPageView *pv = [[CMPageView alloc] init];
    pv.translatesAutoresizingMaskIntoConstraints = NO;
    [self.view addSubview:pv];

    // Layout

    NSArray<NSLayoutConstraint *> *constraints = @[
        [closeBtn.topAnchor constraintEqualToSystemSpacingBelowAnchor:self.view.topAnchor multiplier:2.0],
        [closeBtn.rightAnchor constraintEqualToAnchor:self.view.layoutMarginsGuide.rightAnchor],

        [pv.topAnchor constraintEqualToAnchor:self.view.topAnchor],
        [pv.leftAnchor constraintEqualToAnchor:self.view.leftAnchor],
        [pv.bottomAnchor constraintEqualToAnchor:self.view.bottomAnchor],
        [pv.rightAnchor constraintEqualToAnchor:self.view.rightAnchor],
    ];
    [NSLayoutConstraint activateConstraints:constraints];
}

- (void)setAllowDismissing:(BOOL)allowDismissing {
    _allowDismissing = allowDismissing;
    if (@available(iOS 13.0, *)) {
        self.modalInPresentation = !allowDismissing;
    } else {
        // TODO P0
        // sheets are new in 13 so might be no-op
    }

    _closeButton.hidden = !_allowDismissing;
}

- (void)closeButtonTapped:(UIButton *)sender {
    UIViewController *pvc = self.presentingViewController;
    [pvc dismissViewControllerAnimated:YES completion:nil];
    // TODO dispatch event for dismissed based on sheet name? if so need to get the swipe dismissal as well
}

@end
