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

@end

NS_ASSUME_NONNULL_END
