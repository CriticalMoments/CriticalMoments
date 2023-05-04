#!/bin/sh

# Set working dir
dir="$(cd -P -- "$(dirname -- "$0")" && pwd -P)"
cd $dir

# no op if the go folder hasn't changed. This verifies the code and 
# xcframework build match. 
# Saves time when developing, and ensures edits are built before running.
find . -type f ! -iname ".DS_Store" ! -iname ".go_*_hashlist" -exec md5 {} \; > ./.go_folder_hashlist
diff ./.go_folder_hashlist ./.go_folder_last_build_hashlist > /dev/null 2>&1
folderChanged=$?
if [ $folderChanged -eq 0 ]
then
    echo "No files changed in go directory. New build unnecessary."
    exit 0
fi

echo "Building golang appcore framework..."
# remove prior build and build hashlist
rm -r ./appcore/build/Appcore.xcframework > /dev/null 2>&1
rm ./.go_folder_last_build_hashlist > /dev/null 2>&1

gomobile bind -target ios,iossimulator -iosversion=11 -ldflags '-w -s' -o appcore/build/Appcore.xcframework ./appcore ./cmcore ./cmcore/data_model
buildSuccess=$?

if [ $buildSuccess -eq 0 ]
then
    echo "Build succeeded!"
    find . -type f ! -iname ".DS_Store" ! -iname ".go_*_hashlist" -exec md5 {} \; > ./.go_folder_last_build_hashlist
    exit 0
else
    echo "Build failed. Either fix compiler issues, or re-sync the go folder from remote git repo."
    exit 1
fi

