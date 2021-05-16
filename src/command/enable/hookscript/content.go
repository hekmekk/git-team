package hookscript

import "strings"

const generatedByHeader = "# generated by git-team, do not modify"
const separator = "\n\n"

func withHeader(content string) string {
	return strings.Join([]string{generatedByHeader, strings.Trim(content, "\n")}, separator)
}

var Proxy = withHeader(proxy)

const proxy = `
#!/bin/sh

REAL_LOCAL_HOOK="$(git rev-parse --show-toplevel)/.git/hooks/$(basename ${0})"

if [ -f "${REAL_LOCAL_HOOK}" ]; then
    "${REAL_LOCAL_HOOK}" "${@}" || exit $?
fi

exit 0
`

var PrepareCommitMsg = withHeader(prepareCommitMsg)

const prepareCommitMsg = `
#!/bin/sh

"$(dirname ${0})/prepare-commit-msg-git-team.sh" "${@}" || exit $?

REAL_LOCAL_HOOK="$(git rev-parse --show-toplevel)/.git/hooks/$(basename ${0})"

if [ -f "${REAL_LOCAL_HOOK}" ]; then
   "${REAL_LOCAL_HOOK}" "${@}" || exit $?
fi

exit 0
`

var PrepareCommitMsgGitTeam = withHeader(prepareCommitMsgGitTeam)

const prepareCommitMsgGitTeam = `
#!/bin/sh

activation_scope=$(git config --global team.config.activation-scope)

gitconfig_scope_flag=
if [ "${activation_scope}" = "global" ]; then
        gitconfig_scope_flag=--global
fi

status=$(git config ${gitconfig_scope_flag} team.state.status)

if [ "${status}" != "enabled" ]; then
        exit 0
fi

template=$1
commit_source=$2
commit_hash=$3

# see: https://git-scm.com/docs/githooks#_prepare_commit_msg
# message  - git commit -m|-F
# merge    - git merge (unless ff)
# squash   - git merge --squash
# none     - git commit
# commit   - git commit -c|-C|--amend
# template - git commit -t or if commit.template is set

case "${commit_source}" in
"message" | "merge" | "squash")
        if grep "Co-authored-by:" ${template}; then
                exit 0
        fi

        echo -e -n "\n\n" >> $template
        git config ${gitconfig_scope_flag} --get-all team.state.active-coauthors | while read coauthor; do
                echo -e -n "Co-authored-by: $coauthor\n" >> $template
        done
        ;;
*)
        exit 0
esac

exit 0
`