//
//  CMWeatherPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-10-18.
//

#import <Foundation/Foundation.h>

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMWeatherPropertyProvider : NSObject <CMDynamicPropertyProvider>

- (instancetype)init NS_UNAVAILABLE;
+ (NSDictionary<NSString *, CMWeatherPropertyProvider *> *)allWeatherProviders;

@end

NS_ASSUME_NONNULL_END
