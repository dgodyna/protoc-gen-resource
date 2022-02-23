#!/bin/bash

# This script will be run bazel when building process starts to
# generate key-value information that represents the status of the workspace.
git_rev=$(git describe --exact-match --tags 2>/dev/null || git rev-parse HEAD)
if [[ $? != 0 ]]; then
  exit 1
fi
echo "BUILD_SCM_REVISION ${git_rev}"

short_git_rev=$(git rev-parse --short HEAD)
if [[ $? != 0 ]]; then
  exit 1
fi
echo "BUILD_SHORT_SCM_REVISION ${short_git_rev}"

git_branch=$(git rev-parse --abbrev-ref HEAD)
if [[ $? != 0 ]]; then
  exit 1
fi
echo "GIT_BRANCH ${git_branch}"

buildtime=$(date -u '+%s')
if [[ $? != 0 ]]; then
  exit 1
fi
echo "BUILD_TIME ${buildtime}"

# Check whether there are any uncommitted changes
git diff-index --quiet HEAD --
if [[ $? == 0 ]]; then
  tree_status="Clean"
else
  tree_status="Modified"
fi
echo "BUILD_SCM_STATUS ${tree_status}"

## BUILD_SCM_STAMP should be corrected to be human-readable once we'll release the first stable version
BUILD_SCM_STAMP="${git_rev}"
BUILD_SHORT_SCM_STAMP="${short_git_rev}"

echo "BUILD_SCM_STAMP ${BUILD_SCM_STAMP}"
echo "BUILD_SHORT_SCM_STAMP ${BUILD_SHORT_SCM_STAMP}"
