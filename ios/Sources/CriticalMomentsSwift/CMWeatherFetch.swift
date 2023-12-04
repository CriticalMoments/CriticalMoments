//
//  File.swift
//  
//
//  Created by Steve Cosman on 2023-10-18.
//

#if canImport(WeatherKit)
import WeatherKit
#endif

import CoreLocation

@objc
@objcMembers public class CMWeatherFetch : NSObject {
    private var _current: Any?
    @available(iOS 16.0, *)
    fileprivate var current: CurrentWeather? {
        return _current as? CurrentWeather
    }
    
    @available(iOS 16.0, *)
    public func LoadWeather(location: CLLocation) async -> Bool {
        do {
            let weatherResult = try await WeatherService.shared.weather(for: location, including: .current)
            _current = weatherResult
            return true;
        } catch {
            return false;
        }
    }
    
    @available(iOS 16.0, *)
    public func Temperature() -> NSNumber? {
        let temp = current?.temperature
        if (temp == nil) {
            return nil;
        }
        let celciusVal = temp?.converted(to: UnitTemperature.celsius)
        if (celciusVal != nil) {
            return NSNumber.init(value: celciusVal?.value ?? 0)
        }
        return nil;
    }
    
    @available(iOS 16.0, *)
    public func ApparentTemperature() -> NSNumber? {
        let temp = current?.apparentTemperature
        if (temp == nil) {
            return nil;
        }
        let celciusVal = temp?.converted(to: UnitTemperature.celsius)
        if (celciusVal != nil) {
            return NSNumber.init(value: celciusVal?.value ?? 0)
        }
        return nil;
    }
    
    @available(iOS 16.0, *)
    public func Condition() -> String? {
        return current?.condition.rawValue
    }
    
    @available(iOS 16.0, *)
    public func CloudCover() -> NSNumber? {
        let cloudCover = current?.cloudCover
        if (cloudCover != nil) {
            return NSNumber.init(value: cloudCover!)
        }
        return nil
    }
    
    @available(iOS 16.0, *)
    public func IsDaylight() -> NSNumber? {
        let isDaylight = current?.isDaylight
        if (isDaylight != nil) {
            return NSNumber.init(booleanLiteral: isDaylight!)
        }
        return nil
    }
}
