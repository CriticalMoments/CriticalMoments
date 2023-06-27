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

    UIButton *closeBtn = [UIButton buttonWithType:UIButtonTypeCustom];
    if (@available(iOS 15.0, *)) {
        // don't use UIButtonTypeClose -- might not render on custom bg color. Use secondaryTextColor for contrast.
        UIImageSymbolConfiguration *sc =
            [UIImageSymbolConfiguration configurationWithHierarchicalColor:theme.secondaryTextColor];

        // Size relative to systemFontSize to scale for accessbility.
        sc = [sc configurationByApplyingConfiguration:[UIImageSymbolConfiguration
                                                          configurationWithPointSize:UIFont.systemFontSize * 1.9]];

        UIImage *closeImage = [UIImage systemImageNamed:@"xmark.circle.fill" withConfiguration:sc];

        [closeBtn setImage:closeImage forState:UIControlStateNormal];
    } else {
        // Primary font color here because symbol is visually lighter
        [closeBtn setTitle:@"✕" forState:UIControlStateNormal];
        [closeBtn setTitleColor:theme.primaryTextColor forState:UIControlStateNormal];
        // Size relative to systemFontSize to scale for accessbility.
        closeBtn.titleLabel.font = [UIFont systemFontOfSize:UIFont.systemFontSize * 1.6];
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
        // 44=HIG accessibility recommendation
        [closeBtn.heightAnchor constraintGreaterThanOrEqualToConstant:44],
        [closeBtn.widthAnchor constraintGreaterThanOrEqualToConstant:44],

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
