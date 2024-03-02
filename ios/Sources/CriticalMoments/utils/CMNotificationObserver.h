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

- (instancetype)init NS_UNAVAILABLE;

- (instancetype)initWithCm:(CriticalMoments *)cm;

- (void)start;

@end

NS_ASSUME_NONNULL_END
