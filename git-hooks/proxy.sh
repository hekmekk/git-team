#!/bin/sh

REAL_LOCAL_HOOK="`git rev-parse --show-toplevel`/.git/hooks/`basename ${0}`"

if [ -f "${REAL_LOCAL_HOOK}" ]; then
    "${REAL_LOCAL_HOOK}" "$@"
    result=$?
    if [ $result -ne 0 ]; then
        exit $result
    fi
fi

exit 0
