#!/bin/sh

# Set working dir
dir="$(cd -P -- "$(dirname -- "$0")" && pwd -P)"
cd $dir

pathList=(.go_folder_last_build_hashlist appcore/build/Appcore.xcframework/Info.plist appcore/build/Appcore.xcframework/ios-arm64/Appcore.framework/Versions/A/Appcore appcore/build/Appcore.xcframework/ios-arm64_x86_64-simulator/Appcore.framework/Versions/A/Appcore)

for path in ${pathList[@]}; do
  echo "Adding $path"
  git update-index --no-skip-worktree $path 
  git add $path
  git update-index --skip-worktree $path
done


