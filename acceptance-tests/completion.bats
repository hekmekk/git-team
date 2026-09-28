#!/usr/bin/env bats

# Copyright (C) 2026 Rea Sand
# Licensed under the EUPL

setup() {
  bats_load_library bats-support
  bats_load_library bats-assert
}

@test 'git-team: completion should show the available scripts' {
  run /usr/local/bin/git-team completion

  assert_success
  assert_line --index 6 'COMMANDS:'
  assert_line --index 7 '   bash     Bash completion'
  assert_line --index 8 '   zsh      Zsh completion'
}

@test 'git-team: completion bash should print the bash completion script' {
  run /usr/local/bin/git-team completion bash

  assert_success
  assert_line --index 0 '#!/bin/bash'
  assert_line --index 4 '_git_team() {'
  assert_line --index 17 '}'
  assert_line --index 19 '_git_team_bash_completion() {'
  assert_line --index 32 '}'
  assert_line --index 33 'complete -F _git_team_bash_completion git-team'
}

@test 'git-team: completion zsh should print the zsh completion script' {
  run /usr/local/bin/git-team completion zsh

  assert_success
  assert_line --index 0 '#compdef git-team'
  assert_line --index 3 'function _git-team {'
  assert_line --index 25 '}'
  assert_line --index 26 'compdef _git-team git-team'
  assert_line --index 29 "zstyle ':completion:*:*:git:*' user-commands team:'manage and enhance git commit messages with co-authors'"
}
