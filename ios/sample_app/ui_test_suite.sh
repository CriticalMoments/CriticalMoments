#!/bin/sh

spinner()
{
    local pid=$!

    spin='-\|/'

    i=0
    while kill -0 $pid 2>/dev/null
    do
      i=$(( (i+1) %4 ))
      printf "\r${spin:$i:1}"
      sleep .1
    done

    printf "\b"
}

# iPhone 14 Pro, 16.4
echo "Running UI Tests For iPhone 14 Pro, on iOS 16.4"
xcodebuild -scheme SampleApp -target SampleAppTests -destination 'platform=iOS Simulator,OS=16.4,name=iPhone 14 Pro' '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test &> /dev/null & spinner
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo "Passed\n\n"
else
  echo "Failed\n"
  exit 99
fi

echo "Running UI tests on iPad Pro 12.9, ios 16.4"
xcodebuild -scheme SampleApp -target SampleAppTests -destination 'platform=iOS Simulator,OS=16.4,name=iPad Pro (12.9-inch) (6th generation)' '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test  &> /dev/null & spinner
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo "Passed\n\n"
else
  echo "Failed\n"
  exit 99
fi

echo "Running UI tests on iPhone 11, ios 13.7"
xcodebuild -scheme SampleApp -target SampleAppTests -destination 'platform=iOS Simulator,OS=13.7,name=iPhone 11' '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test  &> /dev/null & spinner
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo "Passed\n\n"
else
  echo "Failed\n"
  exit 99
fi

echo "Running UI tests on iPhone 6s Plus, ios 15.5"
# problems
xcodebuild -scheme SampleApp -target SampleAppTests -destination 'platform=iOS Simulator,OS=15.5,name=iPhone 6s Plus' '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test &> /dev/null & spinner
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo "Passed\n\n"
else
  echo "Failed\n"
  exit 99
fi

echo "Running UI tests on normal iPad (8th gen), ios 14.5"
xcodebuild -scheme SampleApp -target SampleAppTests -destination 'platform=iOS Simulator,OS=14.5,name=iPad (8th generation)' '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test &> /dev/null & spinner
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo "Passed\n\n"
else
  echo "Failed\n"
  exit 99
fi

echo "All Simulator Tests Passed! But you arn't done yet!"
echo " - Manually: Run the snapshot test suite on iPhone 14 Pro HW using xcode"
echo " - Manually: Run manual sanity check on iPhone 6 hardware, running iOS 12. No automation because snapshot test library does not compile to ios 12" 


