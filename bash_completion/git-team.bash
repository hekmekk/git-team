#!/usr/bin/env bash
set +x

#GIT_TEAM_AUTHORS_COMMAND="git authors -a -p"
GIT_TEAM_CURRENT_REV=""
GIT_TEAM_AUTHORS=()

_git_team() {
	COMPREPLY=()

	local cur=${COMP_WORDS[COMP_CWORD]}
	local prev=${COMP_WORDS[COMP_CWORD-1]}
	local aliases="$(/usr/bin/env git config --name-only --get-regexp 'team.alias' | awk -v FS='.' '{ print $3 }')"

	if [[ $prev == \"* ]]; then
		__git_team_coauthor_completion
		return 0
	fi
	for i in $aliases; do
		if [[ $prev == $i ]]; then
			__git_team_coauthor_completion
			return 0
		fi
	done

	case $prev in
		--help)
			COMPREPLY=()
			return 0
			;;
		--version)
			COMPREPLY=()
			return 0
			;;
		disable)
			COMPREPLY=()
			return 0
			;;
		enable)
			__git_team_coauthor_completion
			return 0
			;;
		add)
			COMPREPLY=()
			return 0
			;;
		rm)
			__git_team_coauthor_completion
			return 0
			;;
		list)
			COMPREPLY=()
			return 0
			;;
		*)
			__git_team_coauthor_completion
			local show_flags=true
			if [[ $cur == \"* ]]; then
				show_flags=false
			fi
			for i in $aliases; do
				if [[ $cur == $i ]]; then
					show_flags=false
					break;
				fi
			done
			if [[ $show_flags == true ]]; then
				local flags="add enable disable list rm status --help --version"
				COMPREPLY+=( $(compgen -W "$flags" -- $cur) )
			fi
			return 0
			;;
	esac
}

__git_team_coauthor_completion() {
	if [ $(/usr/bin/env git rev-parse --is-inside-work-tree 2>/dev/null) ]; then
		if [[ $(command -v git-authors) ]]; then
			if [[ $cur == \"* ]]; then
				local rev=$(/usr/bin/env git rev-parse HEAD)
				if [[ $rev != $GIT_TEAM_CURRENT_REV ]]; then
					GIT_TEAM_CURRENT_REV=$rev
					GIT_TEAM_AUTHORS=$(${GIT_TEAM_AUTHORS_COMMAND-git authors -a -p})
					GIT_TEAM_AUTHORS="${GIT_TEAM_AUTHORS//\\ /___}"
				fi

				for author in $GIT_TEAM_AUTHORS; do
					if [[ $author =~ ^$cur ]]; then
						COMPREPLY+=( "${author//___/ }" )
					else
						local cur_space_escaped="${cur// /___}"
						if [[ $author =~ ^$cur_space_escaped ]]; then
							COMPREPLY+=( "${author//___/ }" )
						fi
					fi
				done
				COMPREPLY+=( $(compgen -W "" -- $cur) )
			else
				COMPREPLY=( $(compgen -W "\\\"... ${aliases}" -- $cur) )
			fi
		else
			COMPREPLY=( $(compgen -W "${aliases}" -- $cur) )
		fi
	else
		COMPREPLY=()
	fi
}
