#!/bin/bash

# Script that is used as a GoReleaser's build post-hook to replace non-working
# macOS binaries built on Linux with working ones built on macOS.

set -euo pipefail

BUILD_NAME=$1
BUILD_TARGET=$2

REL_PATH="${BUILD_NAME}_${BUILD_TARGET}/${BUILD_NAME}"
SRC_PATH="$PWD/../macos-binaries/$REL_PATH"
DST_PATH="$PWD/dist/$REL_PATH"

if [[ ${OASIS_CORE_LEDGER_REAL_RELEASE:-false} == true && ${BUILD_TARGET} == darwin_amd64 ]]; then
    echo "Moving binary '$SRC_PATH' to '$DST_PATH'"
    mv $SRC_PATH $DST_PATH
    chmod +x $DST_PATH
fi
