#!/bin/sh

REAL_LOCAL_HOOK="`git rev-parse --show-toplevel`/.git/hooks/`basename ${0}`"

/usr/local/etc/git-team/hooks/prepare-commit-msg-git-team "$@"

if [ -f "${REAL_LOCAL_HOOK}" ]; then
   "${REAL_LOCAL_HOOK}" "$@"
fi

exit 0
