#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

teardown() {
	/usr/local/bin/git-team rm noujz
}

@test "git-team: add should add an assignment to git config" {
	run bash -c "/usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>' &>/dev/null && git config --global team.alias.noujz"
	assert_success
	assert_line 'Mr. Noujz <noujz@mr.se>'
}

@test "git-team: add should create an assigment" {
	run /usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>'
	assert_success
	assert_line "Alias 'noujz' -> 'Mr. Noujz <noujz@mr.se>' has been added."
}

@test "git-team: add should ask for override" {
	/usr/local/bin/git-team add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>' <<< 'y'"
	assert_success
	assert_line "Alias 'noujz' -> 'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [y/N]"
}

@test "git-team: add should fail to create an assigment for an invalidly formatted co-author" {
	run /usr/local/bin/git-team add noujz INVALID-CO-AUTHOR
	assert_failure 255
	assert_line "error: Not a valid coauthor: INVALID-CO-AUTHOR"
}

