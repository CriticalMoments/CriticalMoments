import XCTest
import CoreLocation
@testable import CriticalMomentsSwift

final class CriticalMomentsTests: XCTestCase {
    func testSwiftPing() throws {
        // XCTest Documenation
        // https://developer.apple.com/documentation/xctest

        // Defining Test Cases and Test Methods
        // https://developer.apple.com/documentation/xctest/defining_test_cases_and_test_methods
        let pingResponse = CriticalMoments.swiftPing()
        XCTAssertEqual("swiftPong", pingResponse, "swiftPing method failure, we have basic test issues")
    }
    
    @available(iOS 16.0, *)
    func testWeatherProvider() async throws {
        let ws = CMWeatherFetch()
        let toronto = CLLocation(latitude: 43.651070, longitude:-79.347015)
        let success = await ws.LoadWeather(location: toronto)
        XCTAssert(success, "weather call failed")
        
        let temp = ws.Temperature()
        XCTAssert(temp!.doubleValue > -40.0 && temp!.doubleValue < 50.0, "temp out of range")
        
        let appTemp = ws.ApparentTemperature()
        XCTAssert(appTemp!.doubleValue > -40.0 && appTemp!.doubleValue < 50.0, "temp out of range")
        
        let condition = ws.Condition()
        XCTAssertNotNil(condition)
        XCTAssert(condition!.count > 0)
        
        let cloudCover = ws.CloudCover()
        XCTAssertNotNil(cloudCover)
        XCTAssert(cloudCover!.doubleValue >= 0.0 && cloudCover!.doubleValue <= 1.0)
        
        let isDaylight = ws.IsDaylight()
        XCTAssertNotNil(isDaylight)
    }
    
}
