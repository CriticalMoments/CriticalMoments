//
//  CMPermissionsPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-10-16.
//

#import <Foundation/Foundation.h>

@import AVFoundation;
@import Photos;

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMNotificationPermissionsPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMCapturePermissionsPropertyProvider : NSObject <CMDynamicPropertyProvider>

- (instancetype)init NS_UNAVAILABLE;
- (instancetype)initWithMediaType:(AVMediaType)type;

@end

@interface CMContactsPermissionsPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMPhotosPermissionsPropertyProvider : NSObject <CMDynamicPropertyProvider>

- (instancetype)init NS_DEPRECATED_IOS(2.0, 14.0);
- (instancetype)initWithAccessLevel:(PHAccessLevel)level API_AVAILABLE(ios(14));

@end

NS_ASSUME_NONNULL_END
