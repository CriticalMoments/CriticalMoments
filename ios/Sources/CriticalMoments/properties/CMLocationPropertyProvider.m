//
//  CMLocationPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-10-15.
//

#import "CMLocationPropertyProvider.h"

@import CoreLocation;

@interface CMLocationCache : NSObject <CLLocationManagerDelegate>
@property(nonatomic, strong) CLLocationManager *manager;
@property(nonatomic, strong) dispatch_semaphore_t requestWait;
@property(nonatomic, strong) NSDate *lastErrorTimestamp;
@property(nonatomic, strong) CLPlacemark *reverseGeocodeResponse;
@end

@implementation CMLocationCache

static CMLocationCache *sharedInstance = nil;

+ (CMLocationCache *)shared {
    // avoid lock if we can
    if (sharedInstance) {
        return sharedInstance;
    }

    @synchronized(CMLocationCache.class) {
        if (!sharedInstance) {
            sharedInstance = [[CMLocationCache alloc] init];
        }
        return sharedInstance;
    }
}

- (CLAuthorizationStatus)authStatus {
    if (@available(iOS 14.0, *)) {
        // self.manager can be nil, so use fresh one
        CLLocationManager *manager = [[CLLocationManager alloc] init];
        return manager.authorizationStatus;
    } else {
        return CLLocationManager.authorizationStatus;
    }
}

- (bool)isAuthorized {
    CLAuthorizationStatus as = [self authStatus];
    return as == kCLAuthorizationStatusAuthorizedAlways || as == kCLAuthorizationStatusAuthorizedWhenInUse;
}

- (CLLocation *)getLocationFromCache {
    CLLocation *location = self.manager.location;
    if (!location) {
        return nil;
    }

    // allow for 5 mins of staleness
    NSDate *now = [[NSDate alloc] init];
    if ([location.timestamp compare:[now dateByAddingTimeInterval:-5 * 60]] == NSOrderedDescending) {
        return location;
    }

    return nil;
}

- (CLLocation *)getLocationBlocking {
    // Fail fast if we don't have permissions
    if (![self isAuthorized]) {
        return nil;
    }

    // Try cache
    CLLocation *loc = [self getLocationFromCache];
    if (loc) {
        return loc;
    }

    @synchronized(self) {
        // try cache again, may have populated while waiting on @synchronized
        loc = [self getLocationFromCache];
        if (loc) {
            return loc;
        }

        // Fail fast if we errored in last 9s. Conditions with several properties (lat, long, etc)
        // should not dispatch repeated serial 10s waits when location isn't available (eg, airplane mode).
        NSDate *now = [[NSDate alloc] init];
        if (self.lastErrorTimestamp &&
            [self.lastErrorTimestamp compare:[now dateByAddingTimeInterval:-9]] == NSOrderedDescending) {
            return nil;
        }

        // start a request for location, and wait 10s for it to return
        // Fresh semaphore each time because callback can be called several times, so long-running count
        // will be invalid. We're @synchronized(self) here so safe to set class properties. We only
        // exit the @synchronized(self) after we've been signaled or given up waiting, both of which it's
        // okay for the existing self.requestWait pointer to be replaced
        self.requestWait = dispatch_semaphore_create(0);
        dispatch_async(dispatch_get_main_queue(), ^{
          // https://stackoverflow.com/a/77303256/4076298
          if (!self.manager) {
              self.manager = [[CLLocationManager alloc] init];
              self.manager.delegate = self;
          }
          [self.manager requestLocation];
        });
        dispatch_semaphore_wait(self.requestWait, dispatch_time(DISPATCH_TIME_NOW, 10.0 * NSEC_PER_SEC));
    }

    // May still be nil but at this point we're out of ways to get it
    return [self getLocationFromCache];
}

- (CLPlacemark *)reverseGeocode {
    // try cache before request
    if (self.reverseGeocodeResponse) {
        return self.reverseGeocodeResponse;
    }

    CLLocation *loc = [self getLocationBlocking];
    if (!loc) {
        return nil;
    }

    dispatch_semaphore_t geocodeWait = dispatch_semaphore_create(0);
    CLGeocoder *g = [[CLGeocoder alloc] init];
    [g reverseGeocodeLocation:loc
            completionHandler:^(NSArray<CLPlacemark *> *_Nullable placemarks, NSError *_Nullable error) {
              if (error == nil && placemarks.firstObject) {
                  self.reverseGeocodeResponse = placemarks.firstObject;
              }
              dispatch_semaphore_signal(geocodeWait);
            }];

    dispatch_semaphore_wait(geocodeWait, dispatch_time(DISPATCH_TIME_NOW, 5.0 * NSEC_PER_SEC));
    return self.reverseGeocodeResponse;
}

#pragma mark CLLocationManagerDelegate

- (void)locationManager:(CLLocationManager *)manager didFailWithError:(NSError *)error {
    self.lastErrorTimestamp = [[NSDate alloc] init];
    dispatch_semaphore_signal(self.requestWait);
}

- (void)locationManager:(CLLocationManager *)manager didUpdateLocations:(NSArray<CLLocation *> *)locations {
    // Only signal if the latest location is new enough. Sometimes this is called with a stale location.
    // Should keep waiting for a new enough one if that's the case
    CLLocation *loc = [self getLocationFromCache];
    if (loc) {
        dispatch_semaphore_signal(self.requestWait);
    }
}

@end

@implementation CMLocationPermissionsPropertyProvider

- (BOOL)boolValue {
    return [CMLocationCache.shared isAuthorized];
}

- (long)type {
    return AppcoreLibPropertyProviderTypeBool;
}

@end

@implementation CMLocationPermissionDetailedPropertyProvider

- (NSString *)stringValue {
    CLAuthorizationStatus as = [CMLocationCache.shared authStatus];
    switch (as) {
    case kCLAuthorizationStatusNotDetermined:
        return @"not_determined";
    case kCLAuthorizationStatusRestricted:
        return @"restricted";
    case kCLAuthorizationStatusDenied:
        return @"denied";
    case kCLAuthorizationStatusAuthorizedAlways:
        return @"authorized_always";
    case kCLAuthorizationStatusAuthorizedWhenInUse:
        return @"authorized_when_in_use";
    }
    return @"unknown";
}

- (long)type {
    return AppcoreLibPropertyProviderTypeString;
}

@end

@implementation CMLatitudePropertyProvider

- (NSNumber *)nillableFloatValue {
    CLLocation *loc = [CMLocationCache.shared getLocationBlocking];
    if (!loc) {
        return nil;
    }

    return [NSNumber numberWithDouble:loc.coordinate.latitude];
}

- (long)type {
    return AppcoreLibPropertyProviderTypeFloat;
}

@end

@implementation CMLongitudePropertyProvider

- (NSNumber *)nillableFloatValue {
    CLLocation *loc = [CMLocationCache.shared getLocationBlocking];
    if (!loc) {
        return nil;
    }

    return [NSNumber numberWithDouble:loc.coordinate.longitude];
}

- (long)type {
    return AppcoreLibPropertyProviderTypeFloat;
}

@end

@implementation CMCityPropertyProvider

- (NSString *)stringValue {
    CLPlacemark *place = [CMLocationCache.shared reverseGeocode];
    return place.locality;
}

- (long)type {
    return AppcoreLibPropertyProviderTypeString;
}

@end

@implementation CMRegionPropertyProvider

- (NSString *)stringValue {
    CLPlacemark *place = [CMLocationCache.shared reverseGeocode];
    return place.administrativeArea;
}

- (long)type {
    return AppcoreLibPropertyProviderTypeString;
}

@end

@implementation CMCountryPropertyProvider

- (NSString *)stringValue {
    CLPlacemark *place = [CMLocationCache.shared reverseGeocode];
    return place.ISOcountryCode;
}

- (long)type {
    return AppcoreLibPropertyProviderTypeString;
}

@end
