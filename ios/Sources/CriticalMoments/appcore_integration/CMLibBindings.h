//
//  CMActionDispatcher.h
//
//
//  Created by Steve Cosman on 2023-05-05.
//

#import <Foundation/Foundation.h>

@import Appcore;
#import "CriticalMoments.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMLibBindings : NSObject <AppcoreLibBindings>

/// init is not available. Use initWithCM for all use cases.
- (instancetype)init NS_UNAVAILABLE;

- (instancetype)initWithCM:(CriticalMoments *)cm;

@end

NS_ASSUME_NONNULL_END
