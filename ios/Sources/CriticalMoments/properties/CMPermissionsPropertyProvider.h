//
//  CMPermissionsPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-10-16.
//

#import <Foundation/Foundation.h>

@import AVFoundation;

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMNotificationPermissionsPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMCapturePermissionsPropertyProvider : NSObject <CMDynamicPropertyProvider>

- (instancetype)init NS_UNAVAILABLE;
- (instancetype)initWithMediaType:(AVMediaType)type;

@end

NS_ASSUME_NONNULL_END
