#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

teardown() {
	bash -c "/usr/local/bin/git-team assignments rm noujz || true"
}

@test "git-team: add should add an assignment to git config" {
	run bash -c "/usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>' &>/dev/null && git config --global team.alias.noujz"
	assert_success
	assert_line 'Mr. Noujz <noujz@mr.se>'
}

@test "git-team: add should create an assigment" {
	run /usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>'
	assert_success
	assert_line "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should create an assigment when receiving input from stdin" {
	run bash -c "echo 'noujz Mr. Noujz <noujz@mr.se>' | /usr/local/bin/git-team add"
	assert_success
	assert_line "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should ask for override and apply it if user replies with yes" {
	/usr/local/bin/git-team add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>' <<< yes"
	assert_success
	assert_line "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should force override if the '--force-override' option is used" {
	/usr/local/bin/git-team add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team add --force-override noujz 'Mr. Noujz <noujz@mr.se>'"
	assert_success
	assert_line "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should force override if the '-f' option is used" {
	/usr/local/bin/git-team add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team add -f noujz 'Mr. Noujz <noujz@mr.se>'"
	assert_success
	assert_line "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should ask for override and abort if user replies with no" {
	/usr/local/bin/git-team add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>' <<< no"
	assert_success
	assert_line "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] Nothing changed"
}

@test "git-team: add should ask for override and abort if user replies with anything else" {
	/usr/local/bin/git-team add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>' <<< foo"
	assert_success
	assert_line "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] Nothing changed"
}

@test "git-team: add should ask for override and abort if user just hits ENTER" {
	/usr/local/bin/git-team add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>' <<< ''"
	assert_success
	assert_line "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] Nothing changed"
}

@test "git-team: add should fail to create an assigment for an invalidly formatted co-author" {
	run /usr/local/bin/git-team add noujz INVALID-CO-AUTHOR
	assert_failure 1
	assert_line "error: Not a valid coauthor: INVALID-CO-AUTHOR"
}

