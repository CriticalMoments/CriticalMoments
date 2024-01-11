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
  # Clean optional (usually)
  xcodebuild -scheme SampleApp -target SampleAppTests -destination "$1" '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' clean &> /tmp/critical_moments_clean_log.latest
  xcodebuild -scheme SampleApp -target SampleAppTests -destination "$1" '-only-testing:SampleAppTests/SnapshotTests/testScreenshotAllSampleAppFeatures' test &> /tmp/critical_moments_test_log.latest
  RESULT=$?
  kill $spinnerPid 
  wait $spinnerPid 2>/dev/null
  printf "\r                 \r"
  if [[ "$RECORDING" == "true" ]]; then
    echo "Recording. Continuing\n\n"
  elif [ $RESULT -eq 0 ]; then
    echo "\033[0;32mPassed\033[0m\n"
  else
    # Show error in output
    cat /tmp/critical_moments_test_log.latest
    echo "Issue when running UI tests for $1"
    echo "\033[0;31mFailed\033[0m\nEither fix the bug and re-run, or if change is desired run test suite script after setting up test to record in code and running with RECORDING=true\n"
    exit 99
  fi
}

echo "Pulling snapshot test files from Git LFS."
git lfs pull --include="*" --exclude=""
if $? -ne 0; then
  echo "Error pulling snapshot test files from Git LFS. Exiting. Make sure you have git-lfs installed."
  exit 99
fi

# iPhone 15 Plus, iOS 17.0
runTest 'platform=iOS Simulator,OS=17.0.1,name=iPhone 15 Plus'

# iPhone 14 Pro, 16.4
runTest 'platform=iOS Simulator,OS=16.4,name=iPhone 14 Pro'

# iPad Pro 12.9, ios 16.4
runTest 'platform=iOS Simulator,OS=16.4,name=iPad Pro (12.9-inch) (6th generation)'

# iPhone 11, ios 13.7
echo "Warning: not running on iPhone 11, 13.7. Simulator not supported in xcode 15+"
#runTest 'platform=iOS Simulator,OS=13.7,name=iPhone 11'

# iPhone 6s Plus, ios 15.5
runTest 'platform=iOS Simulator,OS=15.5,name=iPhone 6s Plus'

# normal iPad (8th gen), ios 14.5
echo "Warning: not running on iPad (8th gen), ios 14.5. Simulator not supported in xcode 15+"
#runTest 'platform=iOS Simulator,OS=14.5,name=iPad (8th generation)'

echo "All Simulator Tests Passed! But you arn't done yet!"
echo " - Manually: Run manual sanity check on iPhone 6 hardware, running iOS 12. No automation because snapshot test library does not compile to ios 12" 
