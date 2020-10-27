#!/usr/bin/env bash

# triggered for git team
_git_team() {
	local opts
	opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
	__gitcomp "${opts}"
	compopt +o default
}

# triggered for git-team
_git_team_bash_completion() {
	local cur opts
	COMPREPLY=()
	cur="${COMP_WORDS[COMP_CWORD]}"
	opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
	local IFS=$'\n'
	COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
	return 0
}

complete -F _git_team_bash_completion git-team

