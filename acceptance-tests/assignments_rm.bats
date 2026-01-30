#!/usr/bin/env bats

setup() {
	bats_load_library bats-support
	bats_load_library bats-assert
}

@test "git-team: assignments rm should remove an assigment from git config" {
	/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>'

	run bash -c "/usr/local/bin/git-team assignments rm noujz &>/dev/null && git config --global team.alias.noujz"
	assert_failure 1
	refute_output --regexp '\w+'
}

@test "git-team: assignments rm should remove an assigment" {
	/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>'

	run /usr/local/bin/git-team assignments rm noujz
	assert_success
	assert_line --index 0 "Assignment removed: 'noujz'"
}

@test "git-team: assignments rm should fail for a non-existing alias" {
	run /usr/local/bin/git-team assignments rm noujz
	assert_failure
	assert_line --index 0 "error: no such alias: 'noujz'"
}

