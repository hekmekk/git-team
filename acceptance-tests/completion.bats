#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

@test 'git-team: completion should show the available scripts' {
	run /usr/local/bin/git-team completion

	assert_success
	assert_line --index 6 'COMMANDS:'
	assert_line --index 7 '   bash     Bash completion'
}

@test 'git-team: completion bash should print the bash completion script' {
	run /usr/local/bin/git-team completion bash

	assert_success
	assert_line --index 0 '#!/bin/bash'
	assert_line --index 2 '_git_team() {'
	assert_line --index 17 '_git_team_bash_completion() {'
	assert_line --index 31 'complete -F _git_team_bash_completion git-team'
}

