//
//  CMViewPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-05-24.
//

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMDeviceOrientationPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMInterfaceOrientationPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMDarkModePropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMAppStatePropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMForegroundProvider : NSObject <CMDynamicPropertyProvider>
@end

NS_ASSUME_NONNULL_END
