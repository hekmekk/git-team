#!/bin/sh

REAL_LOCAL_HOOK="`git rev-parse --show-toplevel`/.git/hooks/`basename ${0}`"

if [ -f "${REAL_LOCAL_HOOK}" ]; then
    exec "${REAL_LOCAL_HOOK}" "$@"
fi

exit 0
