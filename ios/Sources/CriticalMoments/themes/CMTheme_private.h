NS_ASSUME_NONNULL_BEGIN

@import Appcore;

@interface CMTheme ()

// _private header prevents exposing these to public SDK.

/**
 :nodoc:
 @param acTheme The appcore theme to convert to CMTheme
 */
+ (CMTheme *)themeFromAppcoreTheme:(DatamodelTheme *)acTheme;

// Theme from config, based on name
/// :nodoc:
+ (CMTheme *)namedThemeFromAppcore:(NSString *)themeName;

@end

NS_ASSUME_NONNULL_END
