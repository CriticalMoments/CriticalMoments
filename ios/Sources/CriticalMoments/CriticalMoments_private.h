//
//  CriticalMoments_private.h
//
//
//  Created by Steve Cosman on 2023-09-29.
//

#import "../CriticalMoments_private.h"
#import "include/CriticalMoments.h"

NS_ASSUME_NONNULL_BEGIN

@import Appcore;

@interface CriticalMoments ()

// _private header prevents exposing these to public SDK.

// Internal only -- use sharedSingleton in product
- (id)initInternal;

// Internal only -- use start in product
- (NSError *)startReturningError;

// Internal only -- for testing and built in event
- (void)sendEvent:(NSString *)eventName
          builtIn:(bool)builtIn
          handler:(void (^_Nullable)(NSError *_Nullable error))handler;

/// :nodoc: access named themes
- (DatamodelTheme *)themeFromConfigByName:(NSString *)name;

// Set the current theme for this CM instance.
// Private, only for internal use (demo app).
/// :nodoc:
- (void)setTheme:(CMTheme *)theme;
// Fetch the current theme for this CM instance
// Private, only for internal use (demo app).
/// :nodoc:
- (CMTheme *)currentTheme;

@end

NS_ASSUME_NONNULL_END
