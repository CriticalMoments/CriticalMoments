//
//  CMTheme.m
//
//
//  Created by Steve Cosman on 2023-04-30.
//

#import "CMTheme.h"

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
    if (currentTheme) {
        return currentTheme;
    }

    @synchronized(CMTheme.class) {
        if (!currentTheme) {
            currentTheme = [[self alloc] init];
        }

        return currentTheme;
    }
}

+ (void)setCurrentTheme:(CMTheme *)theme {
    @synchronized(CMTheme.class) {
        currentTheme = theme;
    }
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

@end