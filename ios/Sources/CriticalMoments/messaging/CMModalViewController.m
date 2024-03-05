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
    pv.buttonCallback = ^(NSString *_Nullable buttonName, int buttonIndex) {
      [weakSelf sendEventsForButtonName:buttonName buttonIndex:buttonIndex];
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

- (void)viewDidDisappear:(BOOL)animated {
    [super viewDidDisappear:animated];

    // Send event for modal closed, only if we're being dismissed (not just for a view unload for backgrounding)
    if (self.beingDismissed || self.isMovingFromParentViewController) {
        if (self.completionEventSender && self.modalName) {
            NSString *closeEventName = [NSString stringWithFormat:@"sub-action:%@:closed", self.modalName];
            [self.completionEventSender sendEvent:closeEventName];
        }
    }
}

- (void)sendEventsForButtonName:(NSString *)buttonName buttonIndex:(int)buttonIndex {
    if (self.completionEventSender && self.modalName) {
        if (buttonName.length > 0) {
            NSString *tapEventName = [NSString stringWithFormat:@"sub-action:%@:button:%@", self.modalName, buttonName];
            [self.completionEventSender sendEvent:tapEventName];
        }
        NSString *tapEventIndexName =
            [NSString stringWithFormat:@"sub-action:%@:button_index:%d", self.modalName, buttonIndex];
        [self.completionEventSender sendEvent:tapEventIndexName];
    }
}

- (void)closeButtonTapped:(UIButton *)sender {
    [self dismissSheet];
}

- (void)dismissSheet {
    if (![NSThread isMainThread]) {
        dispatch_async(dispatch_get_main_queue(), ^{
          [self dismissSheet];
        });
        return;
    }

    UIViewController *pvc = self.presentingViewController;
    [pvc dismissViewControllerAnimated:YES completion:nil];
}

@end
