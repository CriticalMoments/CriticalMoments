//
//  CMBatteryLevelPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-05-22.
//

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMBatteryLevelPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMBatteryStatePropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMLowPowerModePropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

NS_ASSUME_NONNULL_END
