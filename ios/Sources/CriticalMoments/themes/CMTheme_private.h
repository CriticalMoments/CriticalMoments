NS_ASSUME_NONNULL_BEGIN

@import Appcore;

@interface CMTheme ()

// _private header prevents exposing these to public SDK.

+ (CMTheme *)current;

/// Test Theme for e2e testing
+ (CMTheme *)testTheme;

+ (CMTheme *)libaryThemeByName:(NSString *)name;

/**
 Get Theme from Appcore theme
 @param acTheme The appcore theme to convert to CMTheme
 */
+ (CMTheme *)themeFromAppcoreTheme:(DatamodelTheme *)acTheme;

/// Theme from config, based on name
+ (CMTheme *)namedThemeFromAppcore:(NSString *)themeName;

@end

NS_ASSUME_NONNULL_END
