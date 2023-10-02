//
//  CriticalMoments_private.h
//
//
//  Created by Steve Cosman on 2023-09-29.
//

#import "include/CriticalMoments.h"

NS_ASSUME_NONNULL_BEGIN

@import Appcore;

@interface CriticalMoments ()

// _private header prevents exposing these to public SDK.

// Internal only -- use sharedSingleton in product
- (id)initInternal;

// Internal only -- use start in product
- (NSError *)startReturningError;

// Internal only -- for testing
- (void)sendEvent:(NSString *)eventName handler:(void (^)(NSError *_Nullable error))handler;

/// :nodoc: access named themes
- (DatamodelTheme *)themeFromConfigByName:(NSString *)name;

@end

NS_ASSUME_NONNULL_END
