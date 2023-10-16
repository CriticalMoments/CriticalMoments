//
//  CMLocationPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-10-15.
//

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMLocationPermissionsPropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMLocationPermissionDetailedPropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMLatitudePropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMLongitudePropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMCityPropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMRegionPropertyProvider : CMBaseDynamicPropertyProvider
@end

@interface CMCountryPropertyProvider : CMBaseDynamicPropertyProvider
@end

NS_ASSUME_NONNULL_END
