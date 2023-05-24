//
//  CMViewPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-05-24.
//

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMDeviceOrientationPropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMInterfaceOrientationPropertyProvider
    : CMBaseDynamicPropertyProvider
@end

NS_ASSUME_NONNULL_END
