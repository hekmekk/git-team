#!/bin/sh

REAL_LOCAL_HOOK="`git rev-parse --show-toplevel`/.git/hooks/`basename ${0}`"

{{ hooks_dir }}/prepare-commit-msg-git-team "$@"
git_team_hook_result=$?
if [ $git_team_hook_result -ne 0 ]; then
    exit $git_team_hook_result
fi

if [ -f "${REAL_LOCAL_HOOK}" ]; then
   "${REAL_LOCAL_HOOK}" "$@"
    local_hook_result=$?
    if [ $local_hook_result -ne 0 ]; then
        exit $local_hook_result
    fi
fi

exit 0
