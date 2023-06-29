//
//  ThemeDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-30.
//

#import "ThemeDemoScreen.h"
#import "Utils.h"

@import CriticalMoments;
@import UIKit;

@interface ThemeDemoScreen () <UIColorPickerViewControllerDelegate>

@property(nonatomic, copy) void (^currentColorCallback)(UIColor *);

@end

@implementation ThemeDemoScreen

static CMTheme *staticCustomTheme = nil;
+ (CMTheme *)customTheme {
    // avoid lock if we can
    if (staticCustomTheme) {
        return staticCustomTheme;
    }

    @synchronized(ThemeDemoScreen.class) {
        if (!staticCustomTheme) {
            staticCustomTheme = [[CMTheme alloc] init];
        }

        return staticCustomTheme;
    }
}

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Theme Config";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {

    // General

    CMDemoAction *resetThemeAction = [[CMDemoAction alloc] init];
    resetThemeAction.title = @"Reset theme to default";
    resetThemeAction.subtitle = @"Clear all theme changes, restoring default";
    [resetThemeAction addTarget:self action:@selector(resetTheme)];

    CMDemoAction *cannedTheme = [[CMDemoAction alloc] init];
    cannedTheme.title = @"Set demo theme";
    cannedTheme.subtitle = @"Set default theme to a new look.";
    [cannedTheme addTarget:self action:@selector(cannedTheme)];

    CMDemoAction *longBannerAction = [[CMDemoAction alloc] init];
    longBannerAction.title = @"Show banner with current theme";
    longBannerAction.subtitle = @"Display a new banner message to see theme edits.";
    longBannerAction.actionCMActionName = @"long_banner";

    [self addSection:@"General" withActions:@[ resetThemeAction, cannedTheme, longBannerAction ]];

    // Banners

    CMDemoAction *bannerFgColorAction = [[CMDemoAction alloc] init];
    bannerFgColorAction.title = @"Banner foreground color";
    bannerFgColorAction.subtitle = @"Change the banner foreground color";
    [bannerFgColorAction addTarget:self action:@selector(changeBannerFg)];

    CMDemoAction *banneBgColorAction = [[CMDemoAction alloc] init];
    banneBgColorAction.title = @"Banner background color";
    banneBgColorAction.subtitle = @"Change the banner background color";
    [banneBgColorAction addTarget:self action:@selector(changeBannerBg)];

    [self addSection:@"Banner Message Style" withActions:@[ bannerFgColorAction, banneBgColorAction ]];

    CMDemoAction *fontNameAction = [[CMDemoAction alloc] init];
    fontNameAction.title = @"Change font";
    fontNameAction.subtitle = @"Change the font used by UI controls";
    [fontNameAction addTarget:self action:@selector(changeFontName)];

    CMDemoAction *boldFontNameAction = [[CMDemoAction alloc] init];
    boldFontNameAction.title = @"Change bold font";
    boldFontNameAction.subtitle = @"Change the bold font used by UI controls";
    [boldFontNameAction addTarget:self action:@selector(changeBoldFontName)];

    CMDemoAction *fontScaleAction = [[CMDemoAction alloc] init];
    fontScaleAction.title = @"Change font scale";
    fontScaleAction.subtitle = @"Scale the font larger or smaller across all UI";
    [fontScaleAction addTarget:self action:@selector(changeFontScale)];

    [self addSection:@"Fonts" withActions:@[ fontNameAction, boldFontNameAction, fontScaleAction ]];
}

- (void)changeBannerFg {
    UIColor *currentColor = CMTheme.current.bannerForegroundColor;
    [self colorPickerForColor:currentColor
                 withCallback:^(UIColor *color) {
                   CMTheme *customTheme = [ThemeDemoScreen customTheme];
                   customTheme.bannerForegroundColor = color;
                   [CMTheme setCurrentTheme:customTheme];
                   [CMBannerManager.shared removeAllAppWideMessages];
                 }];
}

- (void)changeBannerBg {
    UIColor *currentColor = CMTheme.current.bannerBackgroundColor;
    [self colorPickerForColor:currentColor
                 withCallback:^(UIColor *color) {
                   CMTheme *customTheme = [ThemeDemoScreen customTheme];
                   customTheme.bannerBackgroundColor = color;
                   [CMTheme setCurrentTheme:customTheme];
                   [CMBannerManager.shared removeAllAppWideMessages];
                 }];
}

- (void)colorPickerForColor:(UIColor *)color withCallback:(void (^)(UIColor *))callback {
    if (@available(iOS 14.0, *)) {
        UIColorPickerViewController *colorPicker = [[UIColorPickerViewController alloc] init];
        colorPicker.supportsAlpha = false;
        colorPicker.selectedColor = color;
        self.currentColorCallback = callback;
        colorPicker.delegate = self;
        [Utils.keyWindow.rootViewController presentViewController:colorPicker animated:YES completion:nil];
    } else {
        UIAlertController *alert =
            [UIAlertController alertControllerWithTitle:@"Theme colors demo not available"
                                                message:@"Try this part of the sample app on iOS 14 or newer."
                                         preferredStyle:UIAlertControllerStyleAlert];
        UIAlertAction *defaultAction = [UIAlertAction actionWithTitle:@"OK"
                                                                style:UIAlertActionStyleDefault
                                                              handler:^(UIAlertAction *action){
                                                              }];
        [alert addAction:defaultAction];
        [Utils.keyWindow.rootViewController presentViewController:alert animated:YES completion:nil];
    }
}

- (void)colorPickerViewControllerDidSelectColor:(UIColorPickerViewController *)viewController API_AVAILABLE(ios(14.0)) {
    if (self.currentColorCallback) {
        self.currentColorCallback(viewController.selectedColor);
    }
}

- (void)changeFontName {
    UIAlertController *alert =
        [UIAlertController alertControllerWithTitle:@"Change Font By Name"
                                            message:@"Specify a font by name. See iosfonts.com for "
                                                    @"supported values. An empty resets to default."
                                     preferredStyle:UIAlertControllerStyleAlert];
    [alert addTextFieldWithConfigurationHandler:^(UITextField *_Nonnull textField) {
      textField.placeholder = @"Font name";
      textField.text = @"Baskerville";
    }];

    UIAlertController *__weak weakAlert = alert;
    UIAlertAction *defaultAction = [UIAlertAction actionWithTitle:@"OK"
                                                            style:UIAlertActionStyleDefault
                                                          handler:^(UIAlertAction *action) {
                                                            NSString *newFontName =
                                                                weakAlert.textFields.firstObject.text;
                                                            if (newFontName.length == 0) {
                                                                newFontName = nil;
                                                            }
                                                            CMTheme *customTheme = [ThemeDemoScreen customTheme];
                                                            customTheme.fontName = newFontName;
                                                            [CMTheme setCurrentTheme:customTheme];
                                                            [CMBannerManager.shared removeAllAppWideMessages];
                                                          }];
    [alert addAction:defaultAction];

    [Utils.keyWindow.rootViewController presentViewController:alert animated:YES completion:nil];
}

- (void)changeBoldFontName {
    UIAlertController *alert =
        [UIAlertController alertControllerWithTitle:@"Change Bold Font By Name"
                                            message:@"Specify the 'bold' font by name. See "
                                                    @"iosfonts.com for supported values. An empty "
                                                    @"string resets to default."
                                     preferredStyle:UIAlertControllerStyleAlert];
    [alert addTextFieldWithConfigurationHandler:^(UITextField *_Nonnull textField) {
      textField.placeholder = @"Bold font name";
      textField.text = @"Baskerville-Bold";
    }];

    UIAlertController *__weak weakAlert = alert;
    UIAlertAction *defaultAction = [UIAlertAction actionWithTitle:@"OK"
                                                            style:UIAlertActionStyleDefault
                                                          handler:^(UIAlertAction *action) {
                                                            NSString *newBoldFontName =
                                                                weakAlert.textFields.firstObject.text;
                                                            if (newBoldFontName.length == 0) {
                                                                newBoldFontName = nil;
                                                            }
                                                            CMTheme *customTheme = [ThemeDemoScreen customTheme];
                                                            customTheme.boldFontName = newBoldFontName;
                                                            [CMTheme setCurrentTheme:customTheme];
                                                            [CMBannerManager.shared removeAllAppWideMessages];
                                                          }];
    [alert addAction:defaultAction];

    [Utils.keyWindow.rootViewController presentViewController:alert animated:YES completion:nil];
}

- (void)changeFontScale {
    UIAlertController *alert =
        [UIAlertController alertControllerWithTitle:@"Change font scale"
                                            message:@"Specify a float value to scale UI fonts by "
                                                    @"(example: 0.9 or 1.3). An empty string or "
                                                    @"invalid float return to default (1.0)."
                                     preferredStyle:UIAlertControllerStyleAlert];
    [alert addTextFieldWithConfigurationHandler:^(UITextField *_Nonnull textField) {
      textField.placeholder = @"Font scale factor";
    }];

    UIAlertController *__weak weakAlert = alert;
    UIAlertAction *defaultAction = [UIAlertAction actionWithTitle:@"OK"
                                                            style:UIAlertActionStyleDefault
                                                          handler:^(UIAlertAction *action) {
                                                            NSString *textScale = weakAlert.textFields.firstObject.text;
                                                            float scale = [textScale floatValue];
                                                            if (scale <= 0) {
                                                                scale = 1.0;
                                                            }
                                                            CMTheme *customTheme = [ThemeDemoScreen customTheme];
                                                            customTheme.fontScale = scale;
                                                            [CMTheme setCurrentTheme:customTheme];
                                                            [CMBannerManager.shared removeAllAppWideMessages];
                                                          }];
    [alert addAction:defaultAction];

    [Utils.keyWindow.rootViewController presentViewController:alert animated:YES completion:nil];
}

- (void)resetTheme {
    staticCustomTheme = [[CMTheme alloc] init];
    [CMTheme setCurrentTheme:staticCustomTheme];
    [CMBannerManager.shared removeAllAppWideMessages];
}

- (void)cannedTheme {
    CMTheme *customTheme = [[CMTheme alloc] init];
    customTheme.boldFontName = @"AmericanTypewriter-Bold";
    customTheme.fontName = @"AmericanTypewriter";
    customTheme.fontScale = 1.1;
    customTheme.bannerBackgroundColor = [UIColor blackColor];
    customTheme.bannerForegroundColor = [UIColor whiteColor];
    [CMTheme setCurrentTheme:customTheme];
    [CMBannerManager.shared removeAllAppWideMessages];
}

@end
