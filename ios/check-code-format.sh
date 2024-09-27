#!/usr/bin/env bash

# formats the .m and .h files (limited to ios subdir)

# get root git directory for reference -- allows us to run from other directories
ROOT_DIRECTORY="$(git rev-parse --show-toplevel)"
TARGET_DIRECTORY="$ROOT_DIRECTORY/ios"
cd $TARGET_DIRECTORY

# format
find . -name '*.h' -or -name '*.m' | xargs clang-format --dry-run --Werror -style=file 

FMTSTATUS=$?

if [ $FMTSTATUS -ne 0 ]; then
    echo "clang formatting returned status $FMTSTATUS"
    exit $FMTSTATUS
fi

# json
jq '.' sample_app/SampleApp/starterConfig.json > tmp_file
cmp tmp_file sample_app/SampleApp/starterConfig.json
FMTSTATUS=$?
if [ $FMTSTATUS -ne 0 ]; then
    echo "json formatting returned status $FMTSTATUS"
    exit $FMTSTATUS
fi
jq '.' sample_app/SampleApp/cmDevConfig.json > tmp_file
cmp tmp_file sample_app/SampleApp/cmDevConfig.json
FMTSTATUS=$?
if [ $FMTSTATUS -ne 0 ]; then
    echo "json formatting returned status $FMTSTATUS"
    exit $FMTSTATUS
fi

rm tmp_file

exit 0
