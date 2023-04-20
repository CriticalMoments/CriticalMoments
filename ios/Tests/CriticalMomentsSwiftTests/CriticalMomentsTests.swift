import XCTest
@testable import CriticalMomentsSwift

final class CriticalMomentsTests: XCTestCase {
    func testExample() throws {
        // XCTest Documenation
        // https://developer.apple.com/documentation/xctest

        // Defining Test Cases and Test Methods
        // https://developer.apple.com/documentation/xctest/defining_test_cases_and_test_methods
        let pingResponse = CriticalMoments.swiftPing()
        XCTAssertEqual("swiftPong", pingResponse, "swiftPing method failure, we have basic test issues")
    }
}
