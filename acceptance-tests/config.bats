#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

@test "git-team: config should show the current configuration" {
	run bash -c "git team config"
	assert_success
	assert_line --index 0 'config'
	assert_line --index 1 '─ activation-scope: global'
}

@test "git-team: config activation-scope non-existing-value should fail" {
	run bash -c "git team config activation-scope non-existing-value"
	assert_failure 255
	assert_line "error: Unknown activation-scope 'non-existing-value'"
}

@test "git-team: config activation-scope repo-local should set the activation scope to 'repo-local'" {
	run bash -c "git team config activation-scope repo-local"
	assert_success
	assert_line --index 0 "Configuration updated: 'activation-scope' → 'repo-local'"
}

@test "git-team: config activation-scope repo-local should write the configuration to gitconfig" {
	/usr/local/bin/git-team config activation-scope repo-local

	run bash -c "git config --global --get-regexp team.config | sort"
	assert_success
	assert_line --index 0 'team.config.activation-scope repo-local'
}

@test "git-team: config activation-scope global should set the activation scope to 'global'" {
	run bash -c "git team config activation-scope global"
	assert_success
	assert_line --index 0 "Configuration updated: 'activation-scope' → 'global'"
}

@test "git-team: config activation-scope global should write the configuration to gitconfig" {
	/usr/local/bin/git-team config activation-scope global

	run bash -c "git config --global --get-regexp team.config | sort"
	assert_success
	assert_line --index 0 'team.config.activation-scope global'
}
