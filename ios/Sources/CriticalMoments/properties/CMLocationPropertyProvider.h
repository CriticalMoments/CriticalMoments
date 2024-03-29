//
//  CMLocationPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-10-15.
//

#import "CMBaseDynamicPropertyProvider.h"

@import CoreLocation;

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

@interface CMApproxCityPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMApproxCountryPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMApproxRegionPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMApproxLatitudePropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMApproxLongitudePropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMWeatherPropertyProvider : NSObject <CMDynamicPropertyProvider>

+ (NSDictionary<NSString *, CMWeatherPropertyProvider *> *)allWeatherProviders;
+ (void)setTestLocationOverride:(CLLocation *_Nullable)location;

@end

NS_ASSUME_NONNULL_END
