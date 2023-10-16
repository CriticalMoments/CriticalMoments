//
//  CMLocationPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-10-15.
//

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMLocationPermissionsPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMLocationPermissionDetailedPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMLatitudePropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMLongitudePropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMCityPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMRegionPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMCountryPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

NS_ASSUME_NONNULL_END
