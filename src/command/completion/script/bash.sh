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
              # TODO: need to subsequently suggest the REMAINING aliases or decide against doing so ...
              COMPREPLY=($(compgen -W "$(git team assignments list --only-alias)" --  "${COMP_WORDS[COMP_CWORD]}"))
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
              # TODO: need to subsequently suggest the REMAINING aliases or decide against doing so ...
              COMPREPLY=($(compgen -W "$(git team assignments list --only-alias)" --  "${COMP_WORDS[COMP_CWORD]}"))
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
