#!/bin/bash
# triggered for git team
_git_team() {
  echo "git team: COMP_CWORD=${COMP_CWORD} | COMP_CWORDS=${COMP_WORDS[@]}" >> /tmp/completion.log
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
  echo "git-team: COMP_CWORD=${COMP_CWORD} | COMP_CWORDS=${COMP_WORDS[@]}" >> /tmp/completion.log
  case ${COMP_WORDS[COMP_CWORD-1]} in
    -h | --help)
      COMPREPLY=()
      return
      ;;
  esac
  case $COMP_CWORD in
    1)
      case ${COMP_WORDS[COMP_CWORD]} in
        -*)
          COMPREPLY=($(compgen -W "-h --help" -- "${COMP_WORDS[COMP_CWORD]}"))
          ;;
        *)
          COMPREPLY=($(compgen -W "first1 first2 first3" -- "${COMP_WORDS[COMP_CWORD]}"))
          ;;
      esac
      ;;
    2)
      COMPREPLY=($(compgen -W "second1 second2 second3" -- "${COMP_WORDS[COMP_CWORD]}"))
      ;;
    *)
      COMPREPLY=()
      ;;
  esac
}
complete -F _git_team_bash_completion git-team
