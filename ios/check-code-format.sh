#!/usr/bin/env bash

# formats the .m and .h files (limited to ios subdir)

# get root git directory for reference -- allows us to run from other directories
ROOT_DIRECTORY="$(git rev-parse --show-toplevel)"
TARGET_DIRECTORY="$ROOT_DIRECTORY/ios"
cd $TARGET_DIRECTORY

# format
find . -name '*.h' -or -name '*.m' | xargs clang-format --dry-run --Werror -style=file 

FMTSTATUS=$?

echo "formatting returned status $FMTSTATUS"

exit $FMTSTATUS
