#!/bin/sh

# Set working dir
dir="$(cd -P -- "$(dirname -- "$0")" && pwd -P)"
cd $dir

# no op if the go folder hasn't changed. This verifies the code and 
# xcframework build match. 
# Saves time when developing, and ensures edits are built before running.
mkdir -p ./appcore/build
find . -type f ! -iname ".DS_Store" ! -iname ".go_*_hashlist" -exec md5 {} \; > ./appcore/build/.go_folder_hashlist
diff ./appcore/build/.go_folder_hashlist ./appcore/build/.go_folder_last_build_hashlist > /dev/null 2>&1
folderChanged=$?
if [ $folderChanged -eq 0 ]
then
    echo "No files changed in go directory. New build unnecessary."
    exit 0
fi

echo "Building golang appcore framework..."
# remove prior build and build hashlist
rm -r ./appcore/build/Appcore.xcframework > /dev/null 2>&1
rm ./appcore/build/.go_folder_last_build_hashlist > /dev/null 2>&1

gomobile bind -target ios,iossimulator -iosversion=12 -ldflags '-w -s' -o appcore/build/Appcore.xcframework ./appcore ./cmcore ./cmcore/data_model
buildSuccess=$?

if [ $buildSuccess -ne 0 ]
then
    echo "Build failed. Either fix compiler issues, or re-sync the go folder from remote git repo."
    exit 1
fi

# update the Info.plist to work around https://github.com/golang/go/issues/66018
# this is a temporary fix until the issue is resolved in gomobile
cp ./appcore/build_tools/blankInfo.plist ./appcore/build/Appcore.xcframework/ios-arm64/Appcore.framework/Info.plist
cp ./appcore/build_tools/blankInfo.plist ./appcore/build/Appcore.xcframework/ios-arm64_x86_64-simulator/Appcore.framework/Info.plist 
cp ./appcore/build_tools/exeInfo.plist ./appcore/build/Appcore.xcframework/ios-arm64/Appcore.framework/Resources/Info.plist 
cp ./appcore/build_tools/exeInfo.plist ./appcore/build/Appcore.xcframework/ios-arm64_x86_64-simulator/Appcore.framework/Resources/Info.plist 

echo "Build succeeded!"
find . -type f ! -iname ".DS_Store" ! -iname ".go_*_hashlist" -exec md5 {} \; > ./appcore/build/.go_folder_last_build_hashlist
exit 0
