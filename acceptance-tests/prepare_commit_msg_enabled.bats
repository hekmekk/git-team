#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

setup() {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'
	touch /tmp/COMMIT_MSG
}

teardown() {
	/usr/local/bin/git-team disable
	rm /tmp/COMMIT_MSG
}

@test "prepare-commit-msg: git-team enabled - message" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG message && cat /tmp/COMMIT_MSG"
	assert_success
	assert_line --index 0 'Co-authored-by: A <a@x.y>'
	assert_line --index 1 'Co-authored-by: B <b@x.y>'
	assert_line --index 2 'Co-authored-by: C <c@x.y>'
}

@test "prepare-commit-msg: git-team enabled - none" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team enabled - commit" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG commit && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team enabled - template" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG template && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team enabled - merge" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG merge && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team enabled - squash" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG squash && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

