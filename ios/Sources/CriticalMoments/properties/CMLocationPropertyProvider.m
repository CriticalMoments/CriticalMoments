//
//  CMLocationPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-10-15.
//

#import "CMLocationPropertyProvider.h"
#import "../include/CriticalMoments.h"

@import CoreLocation;
@import CriticalMomentsSwift;

@interface GeoIpPlace : NSObject
@property(nonatomic, strong) NSString *city, *region, *isoCountryCode;
@property(nonatomic, strong) NSNumber *latitude, *longitude;
@end

@implementation GeoIpPlace
@end

@interface CMLocationCache : NSObject <CLLocationManagerDelegate>
@property(nonatomic, strong) CLLocationManager *manager;
@property(nonatomic, strong) dispatch_semaphore_t requestWait;
@property(nonatomic, strong) NSDate *lastErrorTimestamp;
@property(nonatomic, strong) CLPlacemark *reverseGeocodeResponse;

// Approx IP location
@property(nonatomic, strong) NSDate *lastGeoIpErrorTimestamp;
@property(nonatomic, strong) NSDate *lastGeoIpTimestamp;
@property(nonatomic, strong) GeoIpPlace *geoIpPlace;
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

- (GeoIpPlace *)getGeoIpPlaceFromCache {
    GeoIpPlace *ipplace = self.geoIpPlace;
    if (!ipplace) {
        return nil;
    }

    // allow for 20 mins of staleness for geoip
    NSDate *now = [[NSDate alloc] init];
    if ([self.lastGeoIpTimestamp compare:[now dateByAddingTimeInterval:-20 * 60]] == NSOrderedDescending) {
        return ipplace;
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

- (GeoIpPlace *)getApproxLocation {
    GeoIpPlace *ipplace = [self getGeoIpPlaceFromCache];
    if (ipplace) {
        return ipplace;
    }

    @synchronized(self) {
        // try cache again, may have populated while waiting on @synchronized
        ipplace = [self getGeoIpPlaceFromCache];
        if (ipplace) {
            return ipplace;
        }

        // Fail fast if we errored in last 9s. Conditions with several properties (lat, long, etc)
        // should not dispatch repeated serial 10s waits when ip location isn't available (eg, no network).
        NSDate *now = [[NSDate alloc] init];
        if (self.lastGeoIpErrorTimestamp &&
            [self.lastGeoIpErrorTimestamp compare:[now dateByAddingTimeInterval:-9]] == NSOrderedDescending) {
            return nil;
        }

        NSURL *url = [NSURL URLWithString:@"https://api.criticalmoments.io/geo_ip"];
        NSMutableURLRequest *req = [[NSMutableURLRequest alloc] initWithURL:url];
        NSString *apiKey = [CriticalMoments.sharedInstance getApiKey];
        [req setValue:apiKey forHTTPHeaderField:@"X-CM-Api-Key"];

        dispatch_semaphore_t approxSem = dispatch_semaphore_create(0);
        [[[NSURLSession sharedSession] dataTaskWithRequest:req
                                         completionHandler:^(NSData *data, NSURLResponse *response, NSError *reqErr) {
                                           NSHTTPURLResponse *httpResp;
                                           GeoIpPlace *newPlace;
                                           if ([response isKindOfClass:[NSHTTPURLResponse class]]) {
                                               httpResp = (NSHTTPURLResponse *)response;
                                           }
                                           if (!reqErr && data.length > 0 && httpResp && httpResp.statusCode == 200) {
                                               NSError *error = nil;
                                               id jsonObj = [NSJSONSerialization JSONObjectWithData:data
                                                                                            options:0
                                                                                              error:&error];
                                               if (!error && [jsonObj isKindOfClass:[NSDictionary class]]) {
                                                   NSDictionary *jsonResp = (NSDictionary *)jsonObj;
                                                   newPlace = [[GeoIpPlace alloc] init];
                                                   newPlace.city = jsonResp[@"city"];
                                                   newPlace.region = jsonResp[@"region"];
                                                   newPlace.isoCountryCode = jsonResp[@"country"];
                                                   newPlace.latitude = jsonResp[@"latitude"];
                                                   newPlace.longitude = jsonResp[@"longitude"];
                                               }
                                           }

                                           if (newPlace) {
                                               self.geoIpPlace = newPlace;
                                               self.lastGeoIpTimestamp = [[NSDate alloc] init];
                                           } else {
                                               self.lastGeoIpErrorTimestamp = [[NSDate alloc] init];
                                           }

                                           dispatch_semaphore_signal(approxSem);
                                         }] resume];
        dispatch_semaphore_wait(approxSem, dispatch_time(DISPATCH_TIME_NOW, 10.0 * NSEC_PER_SEC));
    }

    // May still be nil but at this point we're out of ways to get it
    return [self getGeoIpPlaceFromCache];
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

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
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

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
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

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeFloat;
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

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeFloat;
}

@end

@implementation CMCityPropertyProvider

- (NSString *)stringValue {
    CLPlacemark *place = [CMLocationCache.shared reverseGeocode];
    return place.locality;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

@end

@implementation CMRegionPropertyProvider

- (NSString *)stringValue {
    CLPlacemark *place = [CMLocationCache.shared reverseGeocode];
    return place.administrativeArea;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

@end

@implementation CMCountryPropertyProvider

- (NSString *)stringValue {
    CLPlacemark *place = [CMLocationCache.shared reverseGeocode];
    return place.ISOcountryCode;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

@end

@implementation CMApproxCityPropertyProvider

- (NSString *)stringValue {
    GeoIpPlace *place = [CMLocationCache.shared getApproxLocation];
    return place.city;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

@end

@implementation CMApproxCountryPropertyProvider

- (NSString *)stringValue {
    GeoIpPlace *place = [CMLocationCache.shared getApproxLocation];
    return place.isoCountryCode;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

@end

@implementation CMApproxRegionPropertyProvider

- (NSString *)stringValue {
    GeoIpPlace *place = [CMLocationCache.shared getApproxLocation];
    return place.region;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeString;
}

@end

@implementation CMApproxLongitudePropertyProvider

- (NSNumber *)nillableFloatValue {
    return [CMLocationCache.shared getApproxLocation].longitude;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeFloat;
}

@end

@implementation CMApproxLatitudePropertyProvider

- (NSNumber *)nillableFloatValue {
    return [CMLocationCache.shared getApproxLocation].latitude;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeFloat;
}

@end

#pragma mark Weather

typedef NS_ENUM(NSInteger, CMWeatherProperty) {
    CMWeatherPropertyCondition,
    CMWeatherPropertyTemperature,
    CMWeatherPropertyApparentTemperature,
    CMWeatherPropertyCloudCover,
    CMWeatherPropertyIsDaylight,
};

API_AVAILABLE(ios(16.0))
@interface CMWeatherCacheItem : NSObject
@property(nonatomic, strong) CMWeatherFetch *weather;
@property(nonatomic, strong) NSDate *date;
@property(nonatomic, strong) CLLocation *location;
@end

@implementation CMWeatherCacheItem
@end

API_AVAILABLE(ios(16.0))
@interface CMWeatherCache : NSObject
@property(nonatomic, strong) NSMutableArray<CMWeatherCacheItem *> *cachedWeather;
@property(nonatomic, strong) NSDate *lastErrorTime;
@end

@implementation CMWeatherCache

API_AVAILABLE(ios(16.0))
static CMWeatherCache *sharedWeatherCache = nil;
+ (CMWeatherCache *)sharedInstance {
    if (sharedWeatherCache) {
        return sharedWeatherCache;
    }
    @synchronized(CMWeatherCache.class) {
        if (!sharedWeatherCache) {
            sharedWeatherCache = [[self alloc] init];
        }
        return sharedWeatherCache;
    }
}

- (instancetype)init {
    self = [super init];
    if (self) {
        self.cachedWeather = [[NSMutableArray alloc] init];
    }
    return self;
}

- (CMWeatherFetch *)getWeatherWithLocaion:(CLLocation *)location {
    if (!location) {
        return nil;
    }

    // one at a time, so we can use cache for next property
    @synchronized(self) {
        // Check the cache
        NSDate *now = [[NSDate alloc] init];
        for (CMWeatherCacheItem *cacheItem in self.cachedWeather) {
            // Cache for 20 mins, but then remove item
            if ([cacheItem.date compare:[now dateByAddingTimeInterval:(-60 * 20)]] == NSOrderedAscending) {
                [self.cachedWeather removeObject:cacheItem];
                continue;
            }

            // Weather cache hit within 1km and 20 mins
            CLLocationDistance distanceInMeters = [location distanceFromLocation:cacheItem.location];
            if (distanceInMeters < 1000) {
                return cacheItem.weather;
            }
        }

        // Fail fast if errored in last 30s. Don't want to make repeated requests if we know network is down
        // or weather service doesn't work.
        if (self.lastErrorTime &&
            [self.lastErrorTime compare:[now dateByAddingTimeInterval:-30]] == NSOrderedDescending) {
            return nil;
        }

        CMWeatherFetch *weather = [self requestWeatherWithLocation:location];
        if (weather) {
            CMWeatherCacheItem *cacheItem = [[CMWeatherCacheItem alloc] init];
            cacheItem.weather = weather;
            cacheItem.date = now;
            cacheItem.location = location;
            [self.cachedWeather addObject:cacheItem];
            return weather;
        }
        return nil;
    }
}

- (CMWeatherFetch *)requestWeatherWithLocation:(CLLocation *)location {
    CMWeatherFetch *weatherFetch = [[CMWeatherFetch alloc] init];
    __block BOOL success = false;
    dispatch_semaphore_t sem = dispatch_semaphore_create(0);
    [weatherFetch LoadWeatherWithLocation:location
                        completionHandler:^(BOOL complete) {
                          success = complete;
                          if (!complete) {
                              self.lastErrorTime = [[NSDate alloc] init]; // now
                          }
                          dispatch_semaphore_signal(sem);
                        }];
    dispatch_semaphore_wait(sem, dispatch_time(DISPATCH_TIME_NOW, 10.0 * NSEC_PER_SEC));

    if (success) {
        return weatherFetch;
    }
    return nil;
}

@end

API_AVAILABLE(ios(16.0))
@interface CMWeatherPropertyProvider ()
@property(nonatomic) CMWeatherProperty property;
@property(nonatomic) BOOL approxAccuracy;
@end

@implementation CMWeatherPropertyProvider

+ (NSDictionary<NSString *, CMWeatherPropertyProvider *> *)allWeatherProviders {
    return @{
        @"weather_temperature" : [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyTemperature
                                                                          forApproxAccuracy:false],
        @"weather_apparent_temperature" :
            [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyApparentTemperature
                                                     forApproxAccuracy:false],
        @"weather_condition" : [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyCondition
                                                                        forApproxAccuracy:false],
        @"weather_cloud_cover" : [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyCloudCover
                                                                          forApproxAccuracy:false],
        @"is_daylight" : [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyIsDaylight
                                                                  forApproxAccuracy:false],
        @"weather_approx_location_temperature" :
            [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyTemperature
                                                     forApproxAccuracy:true],
        @"weather_approx_location_apparent_temperature" :
            [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyApparentTemperature
                                                     forApproxAccuracy:true],
        @"weather_approx_location_condition" :
            [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyCondition
                                                     forApproxAccuracy:true],
        @"weather_approx_location_cloud_cover" :
            [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyCloudCover
                                                     forApproxAccuracy:true],
        @"approx_location_is_daylight" :
            [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyIsDaylight
                                                     forApproxAccuracy:true],
    };
}

static CLLocation *testLocationOverride = nil;
+ (void)setTestLocationOverride:(CLLocation *)location {
    testLocationOverride = location;
}

- (instancetype)initWithWeatherProperty:(CMWeatherProperty)property forApproxAccuracy:(BOOL)approxAccuracy {
    self = [super init];
    if (self) {
        self.property = property;
        self.approxAccuracy = approxAccuracy;
    }
    return self;
}

- (CMWeatherFetch *)getWeather {
    CLLocation *location;
    if (testLocationOverride) {
        location = testLocationOverride;
    } else if (self.approxAccuracy) {
        GeoIpPlace *place = [CMLocationCache.shared getApproxLocation];
        if (place && place.latitude && place.longitude) {
            location = [[CLLocation alloc] initWithLatitude:place.latitude.doubleValue
                                                  longitude:place.longitude.doubleValue];
        }
    } else {
        location = [CMLocationCache.shared getLocationBlocking];
    }

    if (!location) {
        return nil;
    }

    return [CMWeatherCache.sharedInstance getWeatherWithLocaion:location];
}

- (CMPropertyProviderType)type {
    switch (self.property) {
    case CMWeatherPropertyCondition:
        return CMPropertyProviderTypeString;
    case CMWeatherPropertyTemperature:
        return CMPropertyProviderTypeFloat;
    case CMWeatherPropertyApparentTemperature:
        return CMPropertyProviderTypeFloat;
    case CMWeatherPropertyCloudCover:
        return CMPropertyProviderTypeFloat;
    case CMWeatherPropertyIsDaylight:
        return CMPropertyProviderTypeString;
    }
}

- (NSString *)stringValue {
    if (self.property == CMWeatherPropertyCondition) {
        return [self getWeather].Condition;
    } else if (self.property == CMWeatherPropertyIsDaylight) {
        NSNumber *isDaylight = [self getWeather].IsDaylight;
        if (!isDaylight) {
            return @"unknown";
        } else if (isDaylight.boolValue) {
            return @"daylight";
        } else {
            return @"not_daylight";
        }
    }
    return nil;
}

- (NSNumber *)nillableFloatValue {
    switch (self.property) {
    case CMWeatherPropertyTemperature:
        return [self getWeather].Temperature;
    case CMWeatherPropertyApparentTemperature:
        return [self getWeather].ApparentTemperature;
    case CMWeatherPropertyCloudCover:
        return [self getWeather].CloudCover;
    default:
        return nil;
    }
}

@end
