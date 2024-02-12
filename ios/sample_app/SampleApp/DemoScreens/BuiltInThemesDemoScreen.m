//
//  BuiltInThemesDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2024-02-12.
//

#import "BuiltInThemesDemoScreen.h"
#import "Utils.h"

@import CriticalMoments;
@import UIKit;

@implementation BuiltInThemesDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Built In Themes";
        self.infoText = @"After changing the theme, explore other sections of the demo "
                        @"app like modals and banners to see the impact.\n\nBuilt in themes respect system light/dark "
                        @"mode. You can try "
                        @"both varients below, or specify '_light' or '_dark' in config for a sepecific style.";
        self.buttonLink = @"https://docs.criticalmoments.io/actions/themes";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {

    // General

    CMDemoAction *resetThemeAction = [[CMDemoAction alloc] init];
    resetThemeAction.title = @"Reset Theme";
    resetThemeAction.snapshotTitle = @"Reset theme to default";
    resetThemeAction.subtitle = @"Clear all theme changes, restoring default 'system' theme.";
    [resetThemeAction addTarget:self action:@selector(resetThemes)];
    resetThemeAction.skipInUiTesting = true;

    CMDemoAction *longBannerAction = [[CMDemoAction alloc] init];
    longBannerAction.title = @"Show banner";
    longBannerAction.subtitle = @"Display a banner using the current theme.";
    longBannerAction.actionCMActionName = @"short_banner";
    longBannerAction.skipInUiTesting = true;

    [self addSection:@"Reset" withActions:@[ resetThemeAction, longBannerAction ]];

    // Built In Themes
    NSDictionary *themes = @{
        @"elegant" : @"A black and white theme, using a modern font (Avenir Next).",
        @"aqua" : @"A blue/green theme, using default system font.",
        @"jazzy" : @"A purple/pink theme, using default system font.",
        @"sea" : @"A deep blue theme, using default system font.",
        @"honey" : @"A yellow/black theme, using default system font.",
        @"terminal" : @"A theme mirroing the look of a system terminal, using a monospace font.",
        @"system" : @"A theme that matches the system look. Uses the default system font, default system colors, and "
                    @"your app's accent color."
    };

    NSArray<NSString *> *postFixes = @[ @"", @"_dark", @"_light" ];
    NSDictionary *sectionTitles = @{
        @"" : @"Respect System Light/Dark Mode",
        @"_light" : @"Force Light Mode",
        @"_dark" : @"Force Dark Mode",
    };
    for (NSString *postfix in postFixes) {
        NSMutableArray<CMDemoAction *> *builtInThemeActions = [[NSMutableArray alloc] init];

        for (NSString *themeClass in themes) {
            NSString *themeName = [NSString stringWithFormat:@"%@%@", themeClass, postfix];
            NSString *themeDescription = themes[themeClass];

            CMDemoAction *builtInTheme = [[CMDemoAction alloc] init];
            builtInTheme.title = [NSString stringWithFormat:@"'%@' Theme", themeName];
            builtInTheme.snapshotTitle = [NSString stringWithFormat:@"theme_%@", themeName];
            // only test system themes to keep test count down.
            builtInTheme.skipInUiTesting = postfix.length > 0;
            builtInTheme.subtitle = themeDescription;
            builtInTheme.actionBlock = ^{
              [CriticalMoments.sharedInstance setBuiltInTheme:themeName];
              [CriticalMoments.sharedInstance removeAllBanners];

              // Pop a modal so the user can see the theme
              [CriticalMoments.sharedInstance performNamedAction:@"theme_modal" handler:nil];
            };
            [builtInTheme addResetTestTarget:self action:@selector(resetThemes)];
            [builtInThemeActions addObject:builtInTheme];
        }

        [self addSection:sectionTitles[postfix] withActions:builtInThemeActions];
    }
}

- (void)resetThemes {
    // reset theme
    [CriticalMoments.sharedInstance setTheme:[[CMTheme alloc] init]];

    // dismiss the sheets
    [Utils.keyWindow.rootViewController.presentedViewController dismissViewControllerAnimated:NO completion:nil];

    // dismiss the alert
    [Utils.keyWindow.rootViewController dismissViewControllerAnimated:NO completion:nil];

    [CriticalMoments.sharedInstance removeAllBanners];
}

@end
