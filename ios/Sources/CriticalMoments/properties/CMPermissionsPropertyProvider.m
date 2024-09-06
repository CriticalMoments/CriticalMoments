//
//  CMPermissionsPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-10-16.
//

#import "CMPermissionsPropertyProvider.h"

@import UserNotifications;
@import Contacts;
@import CoreBluetooth;
#import <os/log.h>

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

@implementation CMContactsPermissionsPropertyProvider

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

- (NSString *)stringValue {
    CNAuthorizationStatus as = [CNContactStore authorizationStatusForEntityType:CNEntityTypeContacts];
    NSString *result = @"unknown";
    switch (as) {
    case CNAuthorizationStatusNotDetermined:
        result = @"not_determined";
        break;
    case CNAuthorizationStatusRestricted:
        result = @"restricted";
        break;
    case CNAuthorizationStatusDenied:
        result = @"denied";
        break;
    case CNAuthorizationStatusAuthorized:
        result = @"authorized";
        break;
    }
    return result;
}

@end

API_AVAILABLE(ios(14))
@interface CMPhotosPermissionsPropertyProvider ()
@property(nonatomic) PHAccessLevel accessLevel;
@end

@implementation CMPhotosPermissionsPropertyProvider

- (instancetype)init {
    self = [super init];
    return self;
}

- (instancetype)initWithAccessLevel:(PHAccessLevel)level {
    self = [super init];
    if (self) {
        self.accessLevel = level;
    }
    return self;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

- (NSString *)stringValue {
    PHAuthorizationStatus as;
    if (@available(iOS 14.0, *)) {
        PHAccessLevel al = self.accessLevel ? self.accessLevel : PHAccessLevelReadWrite;
        as = [PHPhotoLibrary authorizationStatusForAccessLevel:al];
    } else {
        as = [PHPhotoLibrary authorizationStatus];
    }

    switch (as) {
    case PHAuthorizationStatusNotDetermined:
        return @"not_determined";
    case PHAuthorizationStatusRestricted:
        return @"restricted";
    case PHAuthorizationStatusDenied:
        return @"denied";
    case PHAuthorizationStatusAuthorized:
        return @"authorized";
    case PHAuthorizationStatusLimited:
        return @"limited";
        break;
    }

    return @"unknown";
}

@end

@interface CMCalendarPermissionsPropertyProvider ()
@property(nonatomic) EKEntityType entityType;
@end

@implementation CMCalendarPermissionsPropertyProvider

- (instancetype)initWithEntityType:(EKEntityType)entityType {
    self = [super init];
    if (self) {
        self.entityType = entityType;
    }
    return self;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

- (NSString *)stringValue {
    EKAuthorizationStatus as = [EKEventStore authorizationStatusForEntityType:self.entityType];

    switch (as) {
    case EKAuthorizationStatusNotDetermined:
        return @"not_determined";
    case EKAuthorizationStatusRestricted:
        return @"restricted";
    case EKAuthorizationStatusDenied:
        return @"denied";
    case EKAuthorizationStatusAuthorized: // Same as FullAuthorized, but that constant isn't in 16 sdk
        return @"authorized_full";
    case 4: // EKAuthorizationStatusWriteOnly. allow compiling to iOS 16 SDK. Test added to catch if this ever changes
        return @"authorized_write_only";
    }

    return @"unknown";
}

@end

@implementation CMBluetoothPermissionsPropertyProvider

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

- (NSString *)stringValue {
    if (@available(iOS 13.0, *)) {
        CBManagerAuthorization as;

        if (@available(iOS 13.1, *)) {
            as = CBCentralManager.authorization;
        } else {
            CBManager *m = [[CBCentralManager alloc] init];
            as = m.authorization;
        }

        switch (as) {
        case CBManagerAuthorizationNotDetermined:
            return @"not_determined";
        case CBManagerAuthorizationRestricted:
            return @"restricted";
        case CBManagerAuthorizationDenied:
            return @"denied";
        case CBManagerAuthorizationAllowedAlways:
            return @"authorized";
        }

        return @"unknown";
    } else {
        return @"authorized";
    }
}

#ifdef DEBUG
// Check everything is setup correctly, and log a warning if not.
// Only compiled in debug mode, won't run on release builds.
+ (void)devModeCheckBluetoothSetupCorrectly {
    // Check our 2 IDs are included in the app's Info.plist
    // Don't simply error in callback because it isn't run on simulators, and we want devs to see this.
    NSString *bluetoothUsageDescription =
        [[NSBundle mainBundle] objectForInfoDictionaryKey:@"NSBluetoothAlwaysUsageDescription"];
    if (bluetoothUsageDescription.length == 0) {
        os_log_error(OS_LOG_DEFAULT,
                     "CriticalMoments: Setup Issue\nPlease set NSBluetoothAlwaysUsageDescription in your Info.plist "
                     "file. See docs for details: https://docs.criticalmoments.io/quick-start#update-your-info.plist "
                     "\n\nThis warning log is only in debug builds.");
    }
}
#endif

@end
