#!/bin/bash
# triggered for git team
_git_team() {
  compopt +o default
  case ${COMP_WORDS[COMP_CWORD-1]} in
    -h | --help | -v | --version)
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
          opts=$(compgen -W "-h --help -v --version" -- "${COMP_WORDS[COMP_CWORD]}")
          ;;
        *)
          opts=$(compgen -W "add assignments completion config disable enable list ls remove rm status help h" -- "${COMP_WORDS[COMP_CWORD]}")
          ;;
      esac
      ;;
    3)
      case ${COMP_WORDS[COMP_CWORD-1]} in
        add)
          opts=$(compgen -W "-f --force-override -k --keep-existing" -- "${COMP_WORDS[COMP_CWORD]}")
          ;;
        assignments)
          opts=$(compgen -W "add list ls remove rm help h" -- "${COMP_WORDS[COMP_CWORD]}")
          ;;
        completion)
          opts=$(compgen -W "bash zsh help h" -- "${COMP_WORDS[COMP_CWORD]}")
          ;;
        config)
          opts=$(compgen -W "activation-scope" -- "${COMP_WORDS[COMP_CWORD]}")
          ;;
        disable)
          opts=""
          ;;
        enable)
          case ${COMP_WORDS[COMP_CWORD]} in
            -*)
              opts=$(compgen -W "-A --all" -- "${COMP_WORDS[COMP_CWORD]}")
              ;;
            *)
              opts=$(compgen -W "$(git team assignments list --only-alias)" --  "${COMP_WORDS[COMP_CWORD]}")
              ;;
          esac
          ;;
        list | ls)
          case ${COMP_WORDS[COMP_CWORD]} in
            -*)
              opts=$(compgen -W "-o --only-alias" -- "${COMP_WORDS[COMP_CWORD]}")
              ;;
            *)
              opts=""
              ;;
          esac
          ;;
        remove | rm)
          opts=$(compgen -W "$(git team assignments list --only-alias)" --  "${COMP_WORDS[COMP_CWORD]}")
          ;;
        status)
          opts=""
          ;;
        h | help)
          opts=$(compgen -W "add assignments completion config disable enable list ls remove rm status help h" -- "${COMP_WORDS[COMP_CWORD]}")
          ;;
        *)
          opts=""
          ;;
      esac
      ;;
    4)
      case ${COMP_WORDS[2]} in
        assignments)
          case ${COMP_WORDS[3]} in
            add)
              opts=$(compgen -W "-f --force-override -k --keep-existing" -- "${COMP_WORDS[COMP_CWORD]}")
              ;;
            list | ls)
              case ${COMP_WORDS[COMP_CWORD]} in
                -*)
                  opts=$(compgen -W "-o --only-alias" -- "${COMP_WORDS[COMP_CWORD]}")
                  ;;
                *)
                  opts=""
                  ;;
              esac
              ;;
            remove | rm)
              opts=$(compgen -W "$(git team assignments list --only-alias)" --  "${COMP_WORDS[COMP_CWORD]}")
              ;;
            *)
              opts=""
              ;;
          esac
          ;;
        config)
          case ${COMP_WORDS[3]} in
            activation-scope)
              opts=$(compgen -W "global repo-local" -- "${COMP_WORDS[COMP_CWORD]}")
              ;;
            *)
              opts=""
              ;;
          esac
          ;;
        enable)
          case ${COMP_WORDS[COMP_CWORD]} in
            -*)
              opts=""
              ;;
            *)
              # COMP_WORDS={0:git, 1:team, 2:enable, 3:<alias #1>, ...}
              local used_aliases=("${COMP_WORDS[@]:3:$(( ${#COMP_WORDS[@]} - 4 ))}")
              local remaining_aliases=( $(__determine_remaining_aliases "${used_aliases[@]}") )
              opts=$(compgen -W "$(echo "${remaining_aliases[@]}")" -- "${COMP_WORDS[COMP_CWORD]}")
              ;;
          esac
          ;;
        *)
          opts=""
          ;;
      esac
      ;;
    *)
      case ${COMP_WORDS[2]} in
        enable)
          case ${COMP_WORDS[3]} in
            -*)
              opts=""
              ;;
            *)
              # COMP_WORDS={0:git, 1:team, 2:enable, 3:<alias #1>, ...}
              local used_aliases=("${COMP_WORDS[@]:3:$(( ${#COMP_WORDS[@]} - 4 ))}")
              local remaining_aliases=( $(__determine_remaining_aliases "${used_aliases[@]}") )
              opts=$(compgen -W "$(echo "${remaining_aliases[@]}")" -- "${COMP_WORDS[COMP_CWORD]}")
              ;;
          esac
          ;;
        *)
          opts=""
          ;;
      esac
      ;;
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
              # COMP_WORDS={0:git-team, 1:enable, 2:<alias #1>, ...}
              local used_aliases=("${COMP_WORDS[@]:2:$(( ${#COMP_WORDS[@]} - 3 ))}")
              local remaining_aliases=( $(__determine_remaining_aliases "${used_aliases[@]}") )
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
              # COMP_WORDS={0:git-team, 1:enable, 2:<alias #1>, ...}
              local used_aliases=("${COMP_WORDS[@]:2:$(( ${#COMP_WORDS[@]} - 3 ))}")
              local remaining_aliases=( $(__determine_remaining_aliases "${used_aliases[@]}") )
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


function __determine_remaining_aliases() {
  local used_aliases=("$@")
  local all_aliases=($(git team assignments list --only-alias))
  local remaining_aliases=()

  for alias in "${all_aliases[@]}"; do
    local is_alias_used="no"
    for used_alias in "${used_aliases[@]}"; do
      if [[ "${alias}" = "${used_alias}" ]]; then
        is_alias_used="yes"
        # TODO: why not just break here or smth?
      fi
    done
    if [[ "${is_alias_used}" = "no" ]]; then
      remaining_aliases+=($alias)
    fi
  done

  echo "${remaining_aliases[@]}"
}
