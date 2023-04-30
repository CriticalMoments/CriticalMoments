//
//  CMTheme.h
//  
//
//  Created by Steve Cosman on 2023-04-30.
//

#import <Foundation/Foundation.h>

@import UIKit;

NS_ASSUME_NONNULL_BEGIN

/**
 A class to control the theme/style of our UI controls.
 
 To create a custom theme, create a new instance, modify the properties which control style, and make if your current theme with `setCurrentsTheme`. Swift example of custom theme:

 ```swift
 let customTheme = CMTheme.init()
 customTheme.bannerBackgroundColor = UIColor.red;
 customTheme.bannerForegroundColor = UIColor.white;
 CMTheme.setCurrent(customTheme)
 ```
 */
@interface CMTheme : NSObject

#pragma mark Current Theme

/// Retrieves the current/default theme
+(CMTheme*)current;
/// Sets a theme as current/default
+(void) setCurrentTheme:(CMTheme*)theme;

#pragma mark Banners

/// The background color of banner messages
@property (nonatomic, readwrite) UIColor* bannerBackgroundColor;
/// The foreground color of banner messages
@property (nonatomic, readwrite) UIColor* bannerForegroundColor;

#pragma mark Fonts

/// The default font to use. If nil, this theme will use the system font. Example value: "AvenirNext-Regular". Check out https://iosfontlist.com/ for options, or use your own app-installed font.
/// @warning If you set this, it's highly recommended you also set boldFontName for consistent style.
@property (nonatomic, readwrite) NSString* fontName;

/// The default font to use for bold. If nil, this theme will use the system bold font. Example: "AvenirNext-Bold". Check out https://iosfontlist.com/ for options, or use your own app-installed font.
/// @warning If you set this, it's highly recommended you also set fontName for consistent style.
@property (nonatomic, readwrite) NSString* boldFontName;

/// If YES, will scale the font size based on the user settings in "Settings" app or control center. Defaults to YES. Helps for accessibility.
@property (nonatomic, readwrite) BOOL scaleFontForDynamicType;

/// Scales the font size for all CM UI controls. Defaults to 1.0. Useful if your app's style uses font sizes consistently smaller or larger than the system default sizes.
@property (nonatomic, readwrite) float fontScale;

/// Returns a font respecting the theme settings (fontName, scaleForDynamicType, etc).  Defaults to systemFont.
/// @param fontSize The font size to use. Will scale for dynamic type unless scaleFontForDynamicType is false.
-(UIFont*) fontOfSize:(CGFloat)fontSize;

/// Returns a bold font respecting the theme settings (boldFontName, scaleForDynamicType, etc).  Defaults to boldSystemFont.
/// @param fontSize The font size to use. Will scale for dynamic type unless scaleFontForDynamicType is false.
-(UIFont*) boldFontOfSize:(CGFloat)fontSize;

@end

NS_ASSUME_NONNULL_END
