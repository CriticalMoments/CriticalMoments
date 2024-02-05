//
//  CMNotificationObserver.h
//
//
//  Created by Steve Cosman on 2024-02-05.
//

#import <Foundation/Foundation.h>

#import "../include/CriticalMoments.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMNotificationObserver : NSObject

/// :nodoc:
- (instancetype)init NS_UNAVAILABLE;

/// :nodoc:
- (instancetype)initWithCm:(CriticalMoments *)cm;

/// :nodoc:
- (void)start;

@end

NS_ASSUME_NONNULL_END
