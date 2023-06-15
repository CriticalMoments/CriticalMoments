//
//  CMTheme.m
//
//
//  Created by Steve Cosman on 2023-04-30.
//

#import "CMTheme.h"

#import "../utils/CMUtils.h"

@import Appcore;

@interface CMTheme ()

@property(nonatomic, readwrite) UIColor *primaryColor;

@end

@implementation CMTheme

- (instancetype)init {
    self = [super init];
    if (self) {
        self.scaleFontForDynamicType = YES;
        self.fontScale = 1.0;
    }
    return self;
}

#pragma mark Default Theme

static CMTheme *currentTheme = nil;

+ (CMTheme *)current {
    // avoid lock if we can
    if (!currentTheme) {
        @synchronized(CMTheme.class) {
            if (!currentTheme) {
                currentTheme = [[self alloc] init];
            }
        }
    }

    return [CMTheme themeAdaptedForDarkModeFromTheme:currentTheme];
}

+ (CMTheme *)themeAdaptedForDarkModeFromTheme:(CMTheme *)theme {
    if (!theme.darkModeTheme) {
        return theme;
    }

    if (@available(iOS 12.0, *)) {
        UITraitCollection *tc = UIScreen.mainScreen.traitCollection;
        if (tc.userInterfaceStyle == UIUserInterfaceStyleDark) {
            return theme.darkModeTheme;
        }
    }
    return theme;
}

+ (void)setCurrentTheme:(CMTheme *)theme {
    @synchronized(CMTheme.class) {
        currentTheme = theme;
    }
}

#pragma mark Named Themes From Appcore

+ (CMTheme *)namedThemeFromAppcore:(NSString *)themeName {
    DatamodelTheme *appcoreTheme = [AppcoreSharedAppcore() themeForName:themeName];
    if (appcoreTheme) {
        CMTheme *theme = [CMTheme themeFromAppcoreTheme:appcoreTheme];
        return [CMTheme themeAdaptedForDarkModeFromTheme:theme];
    }
    return nil;
}

#pragma mark Built in themes

+ (CMTheme *)elegantTheme {
    DatamodelTheme *appcoreTheme = DatamodelElegantTheme();
    return [CMTheme themeFromAppcoreTheme:appcoreTheme];
}

+ (CMTheme *)testTheme {
    DatamodelTheme *appcoreTheme = DatamodelTestTheme();
    return [CMTheme themeFromAppcoreTheme:appcoreTheme];
}

#pragma mark Appcore interop

+ (CMTheme *)themeFromAppcoreTheme:(DatamodelTheme *)acTheme {
    CMTheme *theme = [[CMTheme alloc] init];

    // banners
    theme.bannerBackgroundColor = [CMUtils colorFromHexString:acTheme.bannerBackgroundColor];
    theme.bannerForegroundColor = [CMUtils colorFromHexString:acTheme.bannerForegroundColor];

    // fonts
    theme.fontName = acTheme.fontName.length > 0 ? acTheme.fontName : nil;
    theme.boldFontName = acTheme.fontName.length > 0 ? acTheme.boldFontName : nil;
    theme.scaleFontForDynamicType = acTheme.scaleFontForUserPreference;
    theme.fontScale = acTheme.fontScale;

    // dark mode
    if (acTheme.darkModeTheme) {
        theme.darkModeTheme = [CMTheme themeFromAppcoreTheme:acTheme.darkModeTheme];
    }

    return theme;
}

#pragma mark Banners

- (UIColor *)bannerBackgroundColor {
    if (_bannerBackgroundColor) {
        return _bannerBackgroundColor;
    }
    return [UIColor systemYellowColor];
}

- (UIColor *)bannerForegroundColor {
    if (_bannerForegroundColor) {
        return _bannerForegroundColor;
    }
    return [UIColor blackColor];
}

#pragma mark Colors

- (UIColor *)backgroundColor {
    if (_backgroundColor) {
        return _backgroundColor;
    }
    if (@available(iOS 13.0, *)) {
        return [UIColor systemBackgroundColor];
    } else {
        return [UIColor whiteColor];
    }
}

- (UIColor *)primaryTextColor {
    if (_primaryTextColor) {
        return _primaryTextColor;
    }
    if (@available(iOS 13.0, *)) {
        return [UIColor labelColor];
    } else {
        return [UIColor blackColor];
    }
}

- (UIColor *)secondaryTextColor {
    if (_secondaryTextColor) {
        return _secondaryTextColor;
    }
    if (@available(iOS 13.0, *)) {
        return [UIColor secondaryLabelColor];
    } else {
        return [UIColor systemGrayColor];
    }
}

- (void)setPrimaryColor:(UIColor *)color {
    _primaryColor = color;
}

- (UIColor *)primaryColorForView:(UIView *)view {
    if (_primaryColor) {
        return _primaryColor;
    }
    if (@available(iOS 15.0, *)) {
        return [UIColor tintColor];
    } else if (view) {
        return [view tintColor];
    } else {
        return [UIColor systemBlueColor];
    }
}

#pragma mark Fonts

- (UIFont *)fontOfSize:(CGFloat)fontSize {
    UIFont *font;
    if (_fontName) {
        font = [UIFont fontWithName:_fontName size:fontSize];
    }

    if (!font) {
        font = [UIFont systemFontOfSize:fontSize];
    }

    return [self scaleFontForConfig:font];
}

- (UIFont *)boldFontOfSize:(CGFloat)fontSize {
    UIFont *font;
    if (_boldFontName) {
        font = [UIFont fontWithName:_boldFontName size:fontSize];
    }

    if (!font) {
        font = [UIFont boldSystemFontOfSize:fontSize];
    }

    return [self scaleFontForConfig:font];
}

- (UIFont *)scaleFontForConfig:(UIFont *)originalFont {
    UIFont *font = originalFont;

    if (_scaleFontForDynamicType) {
        font = [UIFontMetrics.defaultMetrics scaledFontForFont:font];
    }

    if (_fontScale != 1.0 && _fontScale > 0) {
        font = [font fontWithSize:font.pointSize * _fontScale];
    }

    return font;
}

- (CGFloat)titleFontSize {
    return UIFont.systemFontSize * 2.2;
}

- (CGFloat)subtitleFontSize {
    return UIFont.systemFontSize * 1.4;
}

@end
