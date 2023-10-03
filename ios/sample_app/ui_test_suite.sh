#!/bin/sh

# Exit on first error
set -e

# iPhone 14 Pro, 16.4
xcodebuild -scheme SampleApp -target SampleAppTests -destination 'platform=iOS Simulator,OS=16.4,name=iPhone 14 Pro' '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test

# iPad Pro 12.9, 16.4
xcodebuild -scheme SampleApp -target SampleAppTests -destination 'platform=iOS Simulator,OS=16.4,name=iPad Pro (12.9-inch) (6th generation)' '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test

# iPhone 11, 13.7
xcodebuild -scheme SampleApp -target SampleAppTests -destination 'platform=iOS Simulator,OS=13.7,name=iPhone 11' '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test

# iPhone 6s, 15.5
# Erroring, but in signing. P2
xcodebuild -scheme SampleApp -target SampleAppTests -destination 'platform=iOS Simulator,OS=15.5,name=iPhone 6s' '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test

echo "Manually: Run this suite on iPhone 14 pro HW as well"
echo "Manually: Run manual test on iOS 12 (old iPhone 6 HW) -- snapshot test library does not compile to ios 12"


