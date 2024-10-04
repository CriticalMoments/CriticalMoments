#!/usr/bin/env bash

# formats the .m and .h files (limited to ios subdirectory)

# get root git directory for reference -- allows us to run from other directories
ROOT_DIRECTORY="$(git rev-parse --show-toplevel)"
TARGET_DIRECTORY="$ROOT_DIRECTORY/ios"
cd $TARGET_DIRECTORY

# format
FILE_COUNT="$(find . -name '*.h' -or -name '*.m' | wc -l)"
find . -name '*.h' -or -name '*.m' | xargs clang-format -style=file -i

# format json
jq '.' sample_app/SampleApp/starterConfig.json > tmp_file && mv tmp_file sample_app/SampleApp/starterConfig.json
jq '.' sample_app/SampleApp/cmDevConfig.json > tmp_file && mv tmp_file sample_app/SampleApp/cmDevConfig.json

echo "$FILE_COUNT files were formated with clang-format"
