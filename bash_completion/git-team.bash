#!/usr/bin/env bash
set +x

_git_team() {
	COMPREPLY=()

	local cur=${COMP_WORDS[COMP_CWORD]}
	local prev=${COMP_WORDS[COMP_CWORD-1]}
	local aliases=($(/usr/bin/env git config --global --name-only --get-regexp 'team.alias' | awk -v FS='.' '{ print $3 }'))

	# uncomment for debugging
	# echo "COMP_WORDS: ${COMP_WORDS[@]}" >> /tmp/git-team-completion.debug.log
	# echo "cur: $cur" >> /tmp/git-team-completion.debug.log
	# echo "prev: $prev" >> /tmp/git-team-completion.debug.log

	local remainingAliases=(${aliases[@]})
	for comp_word in ${COMP_WORDS[@]}; do
		for available_alias in ${aliases[@]}; do
			if [[ "$comp_word" == "$available_alias" ]]; then
				remainingAliases=(${remainingAliases[@]/$comp_word})
			fi
		done
	done

	case $cur in
		enable | rm)
			if [[ "${#remainingAliases[@]}" -eq 0 ]]; then
				COMPREPLY=()
				return 0
			fi
			COMPREPLY=("$cur ")
			return 0
			;;
		add | assignments | config | application-scope)
			COMPREPLY=("$cur ")
			return 0
			;;
	esac

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
			COMPREPLY=( $(compgen -W "--all -A ${remainingAliases[*]}" -- $cur) )
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
			return 0
			;;
		rm)
			COMPREPLY=( $(compgen -W "${remainingAliases[*]}" -- $cur) )
			return 0
			;;
		list | ls)
			COMPREPLY=()
			return 0
			;;
		config)
			COMPREPLY+=( $(compgen -W "activation-scope -h" -- $cur) )
			return 0
			;;
		activation-scope)
			COMPREPLY+=( $(compgen -W "global repo-local" -- $cur) )
			return 0
			;;
		*)
			# guard for rm
			if [[ "${COMP_WORDS[COMP_CWORD-2]}" == "rm" ]]; then
				COMPREPLY=()
				return 0
			fi

			local suggestRemainingAliases=false
			for available_alias in ${aliases[@]}; do
				if [[ "$cur" == "$available_alias" || "$prev" == "$available_alias" ]]; then
					suggestRemainingAliases=true
					break
				fi
			done

			if [[ "${suggestRemainingAliases}" == "true" && "${#remainingAliases[@]}" -eq 0 ]]; then
				COMPREPLY=()
				return 0
			fi

			if [[ "${suggestRemainingAliases}" == "true" ]]; then
				matches=($(compgen -W "${aliases[*]}" -- $cur))
				if [[ "${#matches[@]}" -eq 1 ]]; then
					COMPREPLY=("${matches[0]} ")
					return 0
				fi
				# Note: all remainingAliases matching $cur
				COMPREPLY=( $(compgen -W "${remainingAliases[*]}" -- $cur) )
				return 0
			else
				local flags="add assignments enable disable ls rm status config -h --help -v --version"
				COMPREPLY+=( $(compgen -W "${aliases[*]} $flags" -- $cur) )
				return 0
			fi
			;;
	esac
}

complete -F _git_team git-team

