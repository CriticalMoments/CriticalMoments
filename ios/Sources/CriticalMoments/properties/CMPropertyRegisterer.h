//
//  CMDefaultProperties.h
//
//
//  Created by Steve Cosman on 2023-05-20.
//

#import <Foundation/Foundation.h>

@import Appcore;

NS_ASSUME_NONNULL_BEGIN

@interface CMPropertyRegisterer : NSObject

- (instancetype)init NS_UNAVAILABLE;

- (instancetype)initWithAppcore:(AppcoreAppcore *)ac;

- (void)registerDefaultPropertiesToAppcore;

@end

NS_ASSUME_NONNULL_END
