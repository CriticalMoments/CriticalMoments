//
//  CMNetworkingPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-05-24.
//

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMLowDataModePropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMNetworkTypePropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMExpensiveNetworkPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMHasActiveNetworkPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMHasWifiConnectionPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMHasCellConnectionPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

NS_ASSUME_NONNULL_END
