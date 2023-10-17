//
//  CMNetworkingPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-05-24.
//

#import "CMNetworkingPropertyProvider.h"

@import Network;

#define NETWORK_READ_WAIT dispatch_time(DISPATCH_TIME_NOW, 2.0 * NSEC_PER_SEC)

@interface CMNetworkMonitor : NSObject

@property(nonatomic, strong) nw_path_monitor_t monitor;
@property(nonatomic, strong) nw_path_t currentPath;
@property(nonatomic, strong) dispatch_semaphore_t readReadySemaphore;
@property(nonatomic, strong) dispatch_group_t readReadyGroup;

@end

@implementation CMNetworkMonitor

static CMNetworkMonitor *sharedInstance = nil;

+ (CMNetworkMonitor *)shared {
    // avoid lock if we can
    if (sharedInstance) {
        return sharedInstance;
    }

    @synchronized(CMNetworkMonitor.class) {
        if (!sharedInstance) {
            sharedInstance = [[CMNetworkMonitor alloc] init];
        }
        return sharedInstance;
    }
}

- (instancetype)init {
    self = [super init];
    if (self) {
        // Sync code: wait for initial write before allowing any reads
        // using a waitGroup, and semaphor to signal the task in the group
        self.readReadySemaphore = dispatch_semaphore_create(0);
        self.readReadyGroup = dispatch_group_create();
        dispatch_group_async(self.readReadyGroup, dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_BACKGROUND, 0), ^{
          dispatch_semaphore_wait(self.readReadySemaphore, DISPATCH_TIME_FOREVER);
        });

        // Network monitor, which signals the read group once we have data
        // stored
        self.monitor = nw_path_monitor_create();
        nw_path_monitor_set_queue(self.monitor, dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_BACKGROUND, 0));
        __weak CMNetworkMonitor *weakSelf = self;
        nw_path_monitor_set_update_handler(self.monitor, ^(nw_path_t _Nonnull path) {
          weakSelf.currentPath = path;
          dispatch_semaphore_signal(weakSelf.readReadySemaphore);
        });
        nw_path_monitor_start(self.monitor);
    }
    return self;
}

- (void)dealloc {
    nw_path_monitor_cancel(self.monitor);
}

- (bool)hasActiveNetwork {
    bool err = dispatch_group_wait(self.readReadyGroup, NETWORK_READ_WAIT);
    if (err) {
        return NO;
    }
    return nw_path_get_status(self.currentPath) == nw_path_status_satisfied;
}

- (bool)isLowDataMode {
    // low data mode added in ios 13
    if (@available(iOS 13.0, *)) {
        bool err = dispatch_group_wait(self.readReadyGroup, NETWORK_READ_WAIT);
        if (err) {
            return NO;
        }
        return nw_path_is_constrained(self.currentPath);
    } else {
        return NO;
    }
}

- (bool)isExpensiveNetwork {
    bool err = dispatch_group_wait(self.readReadyGroup, NETWORK_READ_WAIT);
    if (err) {
        return NO;
    }
    return nw_path_is_expensive(self.currentPath);
}

- (NSString *)networkType {
    bool err = dispatch_group_wait(self.readReadyGroup, NETWORK_READ_WAIT);
    if (!err) {
        nw_path_t path = self.currentPath;
        if (nw_path_uses_interface_type(path, nw_interface_type_wifi)) {
            return @"wifi";
        } else if (nw_path_uses_interface_type(path, nw_interface_type_cellular)) {
            return @"cellular";
        } else if (nw_path_uses_interface_type(path, nw_interface_type_wired)) {
            return @"wired";
        }
    }
    return @"unknown";
}

- (bool)hasWifiConnection {
    return [self hasConnectionOfType:nw_interface_type_wifi];
}

- (bool)hasCellConnection {
    return [self hasConnectionOfType:nw_interface_type_cellular];
}

- (bool)hasConnectionOfType:(nw_interface_type_t)type {
    bool err = dispatch_group_wait(self.readReadyGroup, NETWORK_READ_WAIT);
    if (err) {
        return NO;
    }
    nw_path_t path = self.currentPath;
    __block bool returnVal = NO;
    nw_path_enumerate_interfaces(path, ^bool(nw_interface_t _Nonnull interface) {
      if (nw_interface_get_type(interface) == type) {
          returnVal = YES;
      }
      // keep iterating until done or found
      return !returnVal;
    });
    return returnVal;
}

@end

@implementation CMLowDataModePropertyProvider

- (BOOL)boolValue {
    return [CMNetworkMonitor.shared isLowDataMode];
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

@end

@implementation CMNetworkTypePropertyProvider

- (NSString *)stringValue {
    return [CMNetworkMonitor.shared networkType];
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

@end

@implementation CMExpensiveNetworkPropertyProvider

- (BOOL)boolValue {
    return [CMNetworkMonitor.shared isExpensiveNetwork];
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

@end

@implementation CMHasActiveNetworkPropertyProvider

- (BOOL)boolValue {
    return [CMNetworkMonitor.shared hasActiveNetwork];
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

@end

@implementation CMHasWifiConnectionPropertyProvider

- (BOOL)boolValue {
    return [CMNetworkMonitor.shared hasWifiConnection];
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

@end

@implementation CMHasCellConnectionPropertyProvider

- (BOOL)boolValue {
    return [CMNetworkMonitor.shared hasCellConnection];
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

@end
