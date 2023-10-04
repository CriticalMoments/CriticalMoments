#!/bin/sh

# Set directory to script path (required)
cd "$(dirname "$0")"

spinner()
{
    spin='▉▊▋▌▍▎▏▎▍▌▋▊▉'
    i=0
    while true
    do
      i=$(( (i+1) %13 ))
      printf "\r${spin:$i:1}${spin:$i:1}${spin:$i:1}${spin:$i:1}${spin:$i:1}${spin:$i:1}${spin:$i:1}${spin:$i:1}"
      sleep .05
    done
}

runTest()
{
  echo "Running UI Tests for $1"
  spinner & spinnerPid=$!
  xcodebuild -scheme SampleApp -target SampleAppTests -destination "$1" '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test &> /tmp/critical_moments_test_log.latest
  RESULT=$?
  kill $spinnerPid 
  wait $spinnerPid 2>/dev/null
  printf "\r                 \r"
  if [ $RESULT -eq 0 ]; then
    echo "\033[0;32mPassed\033[0m\n"
  else
    # Show error in output
    cat /tmp/critical_moments_test_log.latest
    echo "\033[0;31mFailed\033[0m\n"
    exit 99
  fi
}

# iPhone 14 Pro, 16.4
runTest 'platform=iOS Simulator,OS=16.4,name=iPhone 14 Pro'

# iPad Pro 12.9, ios 16.4
runTest 'platform=iOS Simulator,OS=16.4,name=iPad Pro (12.9-inch) (6th generation)'

# iPhone 11, ios 13.7
runTest 'platform=iOS Simulator,OS=13.7,name=iPhone 11'

# iPhone 6s Plus, ios 15.5
runTest 'platform=iOS Simulator,OS=15.5,name=iPhone 6s Plus'

# normal iPad (8th gen), ios 14.5
runTest 'platform=iOS Simulator,OS=14.5,name=iPad (8th generation)'

echo "All Simulator Tests Passed! But you arn't done yet!"
echo " - Manually: Run the snapshot test suite on iPhone 14 Pro HW using xcode"
echo " - Manually: Run manual sanity check on iPhone 6 hardware, running iOS 12. No automation because snapshot test library does not compile to ios 12" 


