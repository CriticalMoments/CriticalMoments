//
//  CMNetworkingPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-05-24.
//

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMLowDataModePropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMNetworkTypePropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMExpensiveNetworkPropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMHasActiveNetworkPropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMHasWifiConnectionPropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMHasCellConnectionPropertyProvider : CMBaseDynamicPropertyProvider
@end

NS_ASSUME_NONNULL_END
