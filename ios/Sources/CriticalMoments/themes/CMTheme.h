//
//  CMTheme.h
//
//
//  Created by Steve Cosman on 2023-04-30.
//

#import <Foundation/Foundation.h>

@import UIKit;

NS_ASSUME_NONNULL_BEGIN

// Max width for iPad, based on readableContentGuide from Apple
#define CM_MAX_TEXT_WIDTH 672

// clang-format off
/**
 A class to control the theme/style of our UI controls.

 See our detailed guide to themes: https://docs.criticalmoments.io/actions/themes
 
 Themes can be created in configuration, and it is not required you create them in code.
 
 To create a custom theme in code, create a new instance, modify the properties which
 control style, and make if your current theme with `setCurrentsTheme`. Swift
 example of custom theme:

 ```swift
 let customTheme = CMTheme.init()
 customTheme.bannerBackgroundColor = UIColor.red;
 customTheme.bannerForegroundColor = UIColor.white;
 CMTheme.setCurrent(customTheme)
 ```
 */
// clang-format on
@interface CMTheme : NSObject

/// :nodoc: Private method. Retrieves the current/default theme.
+ (CMTheme *)current;

#pragma mark Dark Mode

/// A theme to be used if the user has opted to use iOS's dark mode. If unspecified the primary theme will be used in
/// both light and dark modes.
@property(nonatomic, readwrite) CMTheme *_Nullable darkModeTheme;

#pragma mark Colors

/// The background color for views
@property(nonatomic, readwrite) UIColor *backgroundColor;
/// Text color for primary content labels
@property(nonatomic, readwrite) UIColor *primaryTextColor;
/// Text color for secondary content labels
@property(nonatomic, readwrite) UIColor *secondaryTextColor;

/// Fetches a tint color to to apply to buttons, icons, links and more.
/// You must provide the view this will be rendered in, so we can resolve system colors at runtime from view heiarchy if
/// a color is not set explicity in the theme.
- (UIColor *)primaryColorForView:(UIView *)view;
/// Set the primary color for your brand. This is used on buttons, icons, links and more.
/// Should be legible on backgroundColor.
- (void)setPrimaryColor:(UIColor *)color;

/// The background color of banner messages
@property(nonatomic, readwrite) UIColor *bannerBackgroundColor;
/// The foreground color of banner messages
@property(nonatomic, readwrite) UIColor *bannerForegroundColor;

#pragma mark Fonts

/// The default font to use. If nil, this theme will use the system font.
/// Example value: "AvenirNext-Regular". Check out https://iosfontlist.com/ for
/// options, or use your own app-installed font.
/// @warning If you set this, it's highly recommended you also set boldFontName
/// for consistent style.
@property(nonatomic, readwrite) NSString *_Nullable fontName;

/// The default font to use for bold. If nil, this theme will use the system
/// bold font. Example: "AvenirNext-Bold". Check out https://iosfontlist.com/
/// for options, or use your own app-installed font.
/// @warning If you set this, it's highly recommended you also set fontName for
/// consistent style.
@property(nonatomic, readwrite) NSString *_Nullable boldFontName;

/// If YES, will scale the font size based on the user settings in "Settings"
/// app or control center. Defaults to YES. Helps for accessibility.
@property(nonatomic, readwrite) BOOL scaleFontForDynamicType;

/// Scales the font size for all CM UI controls. Defaults to 1.0. Useful if your
/// app's style uses font sizes consistently smaller or larger than the system
/// default sizes.
@property(nonatomic, readwrite) float fontScale;

/// Returns a font respecting the theme settings (fontName, scaleForDynamicType,
/// etc).  Defaults to systemFont.
/// @param fontSize The font size to use. Will scale for dynamic type unless
/// scaleFontForDynamicType is false.
- (UIFont *)fontOfSize:(CGFloat)fontSize;

/// Returns a bold font respecting the theme settings (boldFontName,
/// scaleForDynamicType, etc).  Defaults to boldSystemFont.
/// @param fontSize The font size to use. Will scale for dynamic type unless
/// scaleFontForDynamicType is false.
- (UIFont *)boldFontOfSize:(CGFloat)fontSize;

/// The font size CM uses for titles in pages and sheets
- (CGFloat)titleFontSize;
/// The font size CM uses for body text
- (CGFloat)bodyFontSize;

@end

NS_ASSUME_NONNULL_END
