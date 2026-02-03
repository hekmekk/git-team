#!/bin/sh

if [ -z "${1}" ]; then
  >&2 printf "error: no build version provided\n"
  exit 1
fi

if [ ! `/usr/bin/env git rev-parse --is-inside-work-tree 2>/dev/null` ]; then
  >&2 printf "warn: not inside a git project, build version check skipped\n"
  exit 0
fi

BUILD_VERSION=$1
DERIVED_BUILD_VERSION=$(git describe --tags --abbrev=0 | sed 's/^v//')

if [ "${BUILD_VERSION}" != "${DERIVED_BUILD_VERSION}" ]; then
  >&2 printf "error: provided build version '${BUILD_VERSION}' does not match version derived from git project '${DERIVED_BUILD_VERSION}'\n"
  exit 1
fi

printf "info: build version OK\n"
exit 0
