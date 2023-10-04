#!/bin/bash

goTestCount=$(grep -r "t.Fatal\|t.Error" go | wc -l)
uiTestCount=$(ls -l ./ios/sample_app/SampleAppTests/__Snapshots__/SwiftSnapshotTests/* | wc -l)
iosUnitTestCount=$(grep -r "XCTestExpectation alloc\|XCTAssert\|XCTFail" ios | wc -l)

echo "Test Counts:"
echo "Core Framework Tests: $goTestCount"
echo "SDK UI Tests        : $uiTestCount"
echo "SDK Unit Tests      : $iosUnitTestCount"