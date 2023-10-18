//
//  CMWeatherPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-10-18.
//

#import "CMWeatherPropertyProvider.h"

@import CriticalMomentsSwift;
@import CoreLocation;

// TODO document
// https://developer.apple.com/weatherkit/get-started/
// Entitlement, rate limits.

typedef NS_ENUM(NSInteger, CMWeatherProperty) {
    CMWeatherPropertyCondition,
    CMWeatherPropertyTemperature,
    CMWeatherPropertyApparentTemperature,
    CMWeatherPropertyCloudCover,
    CMWeatherPropertyIsDaylight,
};

@interface CMWeatherCache : NSObject
@property(nonatomic, strong) CMWeatherService *cachedWeather API_AVAILABLE(ios(16.0));
@end

@implementation CMWeatherCache

static CMWeatherCache *sharedInstance = nil;

+ (CMWeatherCache *)sharedInstance {
    // avoid lock if we can
    if (sharedInstance) {
        return sharedInstance;
    }

    @synchronized(CMWeatherCache.class) {
        if (!sharedInstance) {
            sharedInstance = [[self alloc] init];
        }
        return sharedInstance;
    }
}

- (bool)loadWeather {
    if (@available(iOS 16.0, *)) {
        CMWeatherService *ws = [[CMWeatherService alloc] init];
        CLLocation *toronto = [[CLLocation alloc] initWithLatitude:43.651070 longitude:-79.347015];
        [ws LoadWeatherWithLocation:toronto
                  completionHandler:^(BOOL success) {
                    if (success) {
                        NSLog(@"Worked: %lf", [ws Temperature].doubleValue);
                        self.cachedWeather = ws;
                    } else {
                        NSLog(@"failed");
                    }
                  }];
        // TODO wait!
        return true;
    } else {
        return false;
    }
}

@end

@interface CMWeatherPropertyProvider ()
@property(nonatomic) CMWeatherProperty property;
@end

@implementation CMWeatherPropertyProvider

+ (NSDictionary<NSString *, CMWeatherPropertyProvider *> *)allWeatherProviders {
    if (@available(iOS 16.0, *)) {
        return @{
            @"weather_temperature" :
                [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyTemperature],
            @"weather_apparent_temperature" :
                [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyApparentTemperature],
            @"weather_condition" :
                [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyCondition],
            @"weather_cloud_cover" :
                [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyCloudCover],
            @"is_daylight" : [[CMWeatherPropertyProvider alloc] initWithWeatherProperty:CMWeatherPropertyIsDaylight],
        };
    } else {
        return @{};
    }
}

- (instancetype)initWithWeatherProperty:(CMWeatherProperty)property {
    self = [super init];
    if (self) {
        self.property = property;
    }
    return self;
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
        return CMPropertyProviderTypeBool;
    }
}

- (NSString *)stringValue {
    // TODO
    return nil;
}

- (BOOL)boolValue {
    // TODO
    return false;
}

- (NSNumber *)nillableFloatValue {
    // TODO
    [CMWeatherCache.sharedInstance loadWeather];
    return nil;
}

@end
