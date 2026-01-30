#!/usr/bin/env bats

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
	assert_line --index 2 '_git_team() {'
	assert_line --index 15 '}'
	assert_line --index 17 '_git_team_bash_completion() {'
	assert_line --index 30 '}'
	assert_line --index 31 'complete -F _git_team_bash_completion git-team'
}

@test 'git-team: completion zsh should print the zsh completion script' {
	run /usr/local/bin/git-team completion zsh

	assert_success
	assert_line --index 0 '#compdef git-team'
	assert_line --index 1 'function _git-team {'
	assert_line --index 23 '}'
	assert_line --index 24 'compdef _git-team git-team'
	assert_line --index 27 "zstyle ':completion:*:*:git:*' user-commands team:'manage and enhance git commit messages with co-authors'"
}

