//
//  CMPermissionsPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-10-16.
//

#import "CMPermissionsPropertyProvider.h"

@import UserNotifications;

@implementation CMNotificationPermissionsPropertyProvider

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

- (NSString *)stringValue {
    dispatch_semaphore_t sem = dispatch_semaphore_create(0);
    __block NSString *result = @"unknown";
    [UNUserNotificationCenter.currentNotificationCenter
        getNotificationSettingsWithCompletionHandler:^(UNNotificationSettings *_Nonnull settings) {
          UNAuthorizationStatus as = settings.authorizationStatus;
          switch (as) {
          case UNAuthorizationStatusNotDetermined:
              result = @"not_determined";
              break;
          case UNAuthorizationStatusDenied:
              result = @"denied";
              break;
          case UNAuthorizationStatusAuthorized:
              result = @"authorized";
              break;
          case UNAuthorizationStatusProvisional:
              result = @"provisional";
              break;
          case UNAuthorizationStatusEphemeral:
              result = @"ephemeral";
              break;
          }
          dispatch_semaphore_signal(sem);
        }];
    dispatch_semaphore_wait(sem, dispatch_time(DISPATCH_TIME_NOW, 5.0 * NSEC_PER_SEC));
    return result;
}

@end

@interface CMCapturePermissionsPropertyProvider ()
@property(nonatomic) AVMediaType mediaType;
@end

@implementation CMCapturePermissionsPropertyProvider

- (instancetype)initWithMediaType:(AVMediaType)mediaType {
    self = [super init];
    if (self) {
        self.mediaType = mediaType;
    }
    return self;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

- (NSString *)stringValue {
    AVAuthorizationStatus as = [AVCaptureDevice authorizationStatusForMediaType:self.mediaType];
    NSString *result = @"unknown";
    switch (as) {
    case AVAuthorizationStatusNotDetermined:
        result = @"not_determined";
        break;
    case AVAuthorizationStatusRestricted:
        result = @"restricted";
        break;
    case AVAuthorizationStatusDenied:
        result = @"denied";
        break;
    case AVAuthorizationStatusAuthorized:
        result = @"authorized";
        break;
    }
    return result;
}
@end
