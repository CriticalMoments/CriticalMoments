//
//  CriticalMoments.h
//  CriticalMoments
//
//  Created by Steve Cosman on 2023-04-17.
//

#import <Foundation/Foundation.h>

#import "../messaging/CMBannerManager.h"
#import "../messaging/CMBannerMessage.h"
#import "../themes/CMTheme.h"

NS_ASSUME_NONNULL_BEGIN

/// :nodoc:
@interface CriticalMoments : NSObject

// Simple "ping" method for testing end to end integrations
/// :nodoc:
+ (NSString *)objcPing;

// Golang "ping" method for testing end to end integrations
/// :nodoc:
+ (NSString *)goPing;

/**
 Start should be called once you've performed all needed initialization for
 critical moments. Critical moments won't perform actions until it is started.
 This is typically called in AppDelegate didfinishlaunchingwithoptions, but can
 be anywhere you like, as long as the primary root view controler is already
 rendering when you call start.

 Initializtion that should be performed before calling start:

 - Set critical moments API key (required)
 - Set critical moments config URLs (highly recomended)
 - Setup a default theme from code (optional). Can also be done through config
 or not at all.
 */
+ (void)start;

@end

NS_ASSUME_NONNULL_END
