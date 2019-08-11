#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

@test "git-team: remove should remove an assigment from git config" {
	/usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>'

	run bash -c "/usr/local/bin/git-team rm noujz &>/dev/null && git config --global team.alias.noujz"
	assert_failure 1
	refute_output --regexp '\w+'
}

@test "git-team: remove should remove an assigment" {
	/usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>'

	run /usr/local/bin/git-team rm noujz
	assert_success
	assert_line "Alias 'noujz' has been removed."
}

@test "git-team: remove should remove a non-existing assigment without err" {
	run /usr/local/bin/git-team rm noujz
	assert_success
	assert_line "Alias 'noujz' has been removed."
}

