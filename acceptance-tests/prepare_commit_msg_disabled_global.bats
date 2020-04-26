#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

setup() {
	/usr/local/bin/git-team disable
	touch /tmp/COMMIT_MSG
}

teardown() {
	rm /tmp/COMMIT_MSG
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - message" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG message && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - none" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - commit" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG commit && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - template" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG template && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - merge" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG merge && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - squash" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG squash && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

