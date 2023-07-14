#!/bin/bash
# triggered for git team
_git_team() {
  compopt +o default
  case ${COMP_WORDS[COMP_CWORD-1]} in
    -h | --help)
      return
      ;;
  esac
  case $COMP_CWORD in
    1)
      opts=""
      return
      ;;
    2)
      case ${COMP_WORDS[COMP_CWORD]} in
        -*)
          opts=$(compgen -W "-h --help" -- "${COMP_WORDS[COMP_CWORD]}")
          ;;
        *)
          opts=$(compgen -W "first1 first2 first3" -- "${COMP_WORDS[COMP_CWORD]}")
          ;;
      esac
      ;;
    3)
      opts=$(compgen -W "second1 second2 second3" -- "${COMP_WORDS[COMP_CWORD]}")
      ;;
    *)
      opts=""
  esac
  __gitcomp "${opts}"
}
# triggered for git-team
_git_team_bash_completion() {
  case ${COMP_WORDS[COMP_CWORD-1]} in
    -h | --help | -v | --version)
      COMPREPLY=()
      return
      ;;
  esac
  case $COMP_CWORD in
    1)
      case ${COMP_WORDS[COMP_CWORD]} in
        -*)
          COMPREPLY=($(compgen -W "-h --help -v --version" -- "${COMP_WORDS[COMP_CWORD]}"))
          ;;
        *)
          COMPREPLY=($(compgen -W "add assignments completion config disable enable list ls remove rm status help h" -- "${COMP_WORDS[COMP_CWORD]}"))
          ;;
      esac
      ;;
    2)
      case ${COMP_WORDS[COMP_CWORD-1]} in
        add)
          COMPREPLY=($(compgen -W "-f --force-override -k --keep-existing" -- "${COMP_WORDS[COMP_CWORD]}"))
          ;;
        assignments)
          COMPREPLY=($(compgen -W "add list ls remove rm help h" -- "${COMP_WORDS[COMP_CWORD]}"))
          ;;
        completion)
          COMPREPLY=($(compgen -W "bash zsh help h" -- "${COMP_WORDS[COMP_CWORD]}"))
          ;;
        config)
          COMPREPLY=($(compgen -W "activation-scope" -- "${COMP_WORDS[COMP_CWORD]}"))
          ;;
        disable)
          COMPREPLY=()
          ;;
        enable)
          case ${COMP_WORDS[COMP_CWORD]} in
            -*)
              COMPREPLY=($(compgen -W "-A --all" -- "${COMP_WORDS[COMP_CWORD]}"))
              ;;
            *)
              COMPREPLY=($(compgen -W "$(git team assignments list --only-alias)" --  "${COMP_WORDS[COMP_CWORD]}"))
              ;;
          esac
          ;;
        list | ls)
          case ${COMP_WORDS[COMP_CWORD]} in
            -*)
              COMPREPLY=($(compgen -W "-o --only-alias" -- "${COMP_WORDS[COMP_CWORD]}"))
              ;;
            *)
              COMPREPLY=()
              ;;
          esac
          ;;
        remove | rm)
          COMPREPLY=($(compgen -W "$(git team assignments list --only-alias)" --  "${COMP_WORDS[COMP_CWORD]}"))
          ;;
        status)
          COMPREPLY=()
          ;;
        h | help)
          COMPREPLY=($(compgen -W "add assignments completion config disable enable list ls remove rm status help h" -- "${COMP_WORDS[COMP_CWORD]}"))
          ;;
        *)
          COMPREPLY=()
          ;;
      esac
      ;;
    3)
      case ${COMP_WORDS[1]} in
        assignments)
          case ${COMP_WORDS[2]} in
            add)
              COMPREPLY=($(compgen -W "-f --force-override -k --keep-existing" -- "${COMP_WORDS[COMP_CWORD]}"))
              ;;
            list | ls)
              case ${COMP_WORDS[COMP_CWORD]} in
                -*)
                  COMPREPLY=($(compgen -W "-o --only-alias" -- "${COMP_WORDS[COMP_CWORD]}"))
                  ;;
                *)
                  COMPREPLY=()
                  ;;
              esac
              ;;
            remove | rm)
              COMPREPLY=($(compgen -W "$(git team assignments list --only-alias)" --  "${COMP_WORDS[COMP_CWORD]}"))
              ;;
            *)
              COMPREPLY=()
              ;;
          esac
          ;;
        config)
          case ${COMP_WORDS[2]} in
            activation-scope)
              COMPREPLY=($(compgen -W "global repo-local" -- "${COMP_WORDS[COMP_CWORD]}"))
              ;;
            *)
              COMPREPLY=()
              ;;
          esac
          ;;
        enable)
          case ${COMP_WORDS[COMP_CWORD]} in
            -*)
              COMPREPLY=()
              ;;
            *)
              # note: we're here for git-team enable <ALIAS #1>
              FROM_INDEX=2
              TO_INDEX=$(( ${#COMP_WORDS[@]} - 3 ))
              USED_ALIASES=("${COMP_WORDS[@]:${FROM_INDEX}:${TO_INDEX}}")
              local all_aliases=($(git team assignments list --only-alias))
              local remaining_aliases=()
              for alias in "${all_aliases[@]}"; do
                local is_alias_used="no"
                for used_alias in "${USED_ALIASES[@]}"; do
                  if [[ "${alias}" = "${used_alias}" ]]; then
                    is_alias_used="yes"
                    # TODO: why not just break here or smth?
                  fi
                done
                if [[ "${is_alias_used}" = "no" ]]; then
                  remaining_aliases+=($alias)
                fi
              done
              COMPREPLY=($(compgen -W "$(echo "${remaining_aliases[@]}")" -- "${COMP_WORDS[COMP_CWORD]}"))
              ;;
          esac
          ;;
        *)
          COMPREPLY=()
          ;;
      esac
      ;;
    *)
      case ${COMP_WORDS[1]} in
        enable)
          case ${COMP_WORDS[2]} in
            -*)
              COMPREPLY=()
              ;;
            *)
              # note: we're here for git-team enable <ALIAS #1> <ALIAS #n>
              FROM_INDEX=2
              TO_INDEX=$(( ${#COMP_WORDS[@]} - 3 ))
              USED_ALIASES=("${COMP_WORDS[@]:${FROM_INDEX}:${TO_INDEX}}")
              local all_aliases=($(git team assignments list --only-alias))
              local remaining_aliases=()
              for alias in "${all_aliases[@]}"; do
                local is_alias_used="no"
                for used_alias in "${USED_ALIASES[@]}"; do
                  if [[ "${alias}" = "${used_alias}" ]]; then
                    is_alias_used="yes"
                    # TODO: why not just break here or smth?
                  fi
                done
                if [[ "${is_alias_used}" = "no" ]]; then
                  remaining_aliases+=($alias)
                fi
              done
              COMPREPLY=($(compgen -W "$(echo "${remaining_aliases[@]}")" -- "${COMP_WORDS[COMP_CWORD]}"))
              ;;
          esac
          ;;
        *)
          COMPREPLY=()
          ;;
      esac
      ;;
  esac
}
complete -F _git_team_bash_completion git-team
