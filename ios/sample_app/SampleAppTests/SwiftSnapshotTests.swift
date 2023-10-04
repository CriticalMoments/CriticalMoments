//
//  SwiftSnapshotTests.swift
//  SampleAppTests
//
//  Created by Steve Cosman on 2023-10-02.
//

import Foundation
import SnapshotTesting
import XCTest

// Small util class making the awesome SnapshotTesting library easy to use from obj-c

@objc
class CMSnapshotWrapper: XCTestCase {
        
    @objc
    func assertSnapshotImage(
        of value: UIImage,
        named: String
    ) {
        let failure = verifySnapshot(
          of: value,
          as: .image(precision: 0.995),
          named: named,
          record: false, // Modify to record new tests/fixes
          timeout: 5,
          testName: "demolith"
        )
        guard let message = failure else { return }
        XCTFail(message)
    }
}
