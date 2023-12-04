//
//  File.swift
//  
//
//  Created by Steve Cosman on 2023-10-18.
//

import Foundation

import WeatherKit
import CoreLocation

@available(iOS 16.0, *)
@objc
@objcMembers public class CMWeatherFetch : NSObject {
    var current: CurrentWeather?
    
    public func LoadWeather(location: CLLocation) async -> Bool {
        do {
            let weatherResult = try await WeatherService.shared.weather(for: location, including: .current)
            current = weatherResult
            return true;
        } catch {
            return false;
        }
    }
    
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
    
    public func Condition() -> String? {
        return current?.condition.rawValue
    }
    
    public func CloudCover() -> NSNumber? {
        let cloudCover = current?.cloudCover
        if (cloudCover != nil) {
            return NSNumber.init(value: cloudCover!)
        }
        return nil
    }
    
    public func IsDaylight() -> NSNumber? {
        let isDaylight = current?.isDaylight
        if (isDaylight != nil) {
            return NSNumber.init(booleanLiteral: isDaylight!)
        }
        return nil
    }
}
