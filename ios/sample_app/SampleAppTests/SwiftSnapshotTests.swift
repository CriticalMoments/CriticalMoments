//
//  SwiftSnapshotTests.swift
//  SampleAppTests
//
//  Created by Steve Cosman on 2023-10-02.
//

import Foundation
import SnapshotTesting
import XCTest

@objc
class CMSnapshotWrapper: XCTestCase {
  func testMyViewController() {
    //let vc = MyViewController()

      let vc = UIViewController()
    assertSnapshot(of: vc, as: .image)
  }
    
    @objc
    func assertSnapshot(
        of value: UIWindow,
        named: String
    ) {
        print("test")
        let failure = verifySnapshot(
          of: value,
          as: .image,
          named: named
        )
        guard let message = failure else { return }
        XCTFail(message)
    }
    
    @objc
    func assertSnapshotImage(
        of value: UIImage,
        named: String
    ) {
        print("test")
        let failure = verifySnapshot(
          of: value,
          as: .image,
          named: named
        )
        guard let message = failure else { return }
        XCTFail(message)
    }
    
    @objc
    func assertSnapshotVC(
        of value: UIViewController,
        named: String
    ) {
        print("test")
        let failure = verifySnapshot(
          of: value,
          as: .image,
          named: named
        )
        guard let message = failure else { return }
        XCTFail(message)
    }
    
    public func assertSnapshot<Value, Format>(
      of value: @autoclosure () throws -> Value,
      as snapshotting: Snapshotting<Value, Format>,
      named name: String? = nil,
      record recording: Bool = false,
      timeout: TimeInterval = 5,
      file: StaticString = #file,
      testName: String = #function,
      line: UInt = #line
    ) {
      let failure = verifySnapshot(
        of: try value(),
        as: snapshotting,
        named: name,
        record: recording,
        timeout: timeout,
        file: file,
        testName: testName,
        line: line
      )
      guard let message = failure else { return }
      XCTFail(message, file: file, line: line)
    }
}
