#!/usr/bin/env bash
set +x

GIT_TEAM_CURRENT_REV=""
GIT_TEAM_AUTHORS=()

_git_team() {
	COMPREPLY=()

	local cur=${COMP_WORDS[COMP_CWORD]}
	local prev=${COMP_WORDS[COMP_CWORD-1]}
	local aliases="$(/usr/bin/env git config --global --name-only --get-regexp 'team.alias' | awk -v FS='.' '{ print $3 }')"

	for i in $aliases; do
		if [[ $prev == $i ]]; then
			COMPREPLY=( $(compgen -W "${aliases}" -- $cur) )
			return 0
		fi
	done

	case $prev in
		-h | --help)
			COMPREPLY=()
			return 0
			;;
		-v | --version)
			COMPREPLY=()
			return 0
			;;
		disable)
			COMPREPLY=()
			return 0
			;;
		enable)
			COMPREPLY=( $(compgen -W "${aliases}" -- $cur) )
			return 0
			;;
		add)
			COMPREPLY=( $(compgen -W "--force-override -f" -- $cur) )
			return 0
			;;
		-f | --force-override)
			COMPREPLY=()
			return 0
			;;
		assignments)
			COMPREPLY+=( $(compgen -W "add rm ls" -- $cur) )
			;;
		rm)
			COMPREPLY=( $(compgen -W "${aliases}" -- $cur) )
			return 0
			;;
		list | ls)
			COMPREPLY=()
			return 0
			;;
		*)
			COMPREPLY=( $(compgen -W "${aliases}" -- $cur) )
			local show_flags=true
			for i in $aliases; do
				if [[ $cur == $i ]]; then
					show_flags=false
					break;
				fi
			done
			if [[ $show_flags == true ]]; then
				local flags="add assignments enable disable ls rm status -h --help -v --version"
				COMPREPLY+=( $(compgen -W "$flags" -- $cur) )
			fi
			return 0
			;;
	esac
}

