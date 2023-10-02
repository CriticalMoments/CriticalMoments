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

// Internal tests only
- (id)initInternal;
- (NSError *)startReturningError;

/// :nodoc: access named themes
- (DatamodelTheme *)themeFromConfigByName:(NSString *)name;

@end

NS_ASSUME_NONNULL_END