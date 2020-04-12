#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

@test "git-team: config should show the current configuration" {
	run bash -c "git team config"
	assert_success
	assert_line --index 0 'config'
	assert_line --index 1 '─ activation-scope: global'
}

