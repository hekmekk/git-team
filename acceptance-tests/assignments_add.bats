#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

teardown() {
	bash -c "/usr/local/bin/git-team assignments rm noujz || true"
}

@test "git-team: assignments add should add an assignment to git config" {
	run bash -c "/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>' &>/dev/null && git config --global team.alias.noujz"
	assert_success
	assert_line 'Mr. Noujz <noujz@mr.se>'
}

@test "git-team: assignments add should create a new assigment" {
	run /usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>'
	assert_success
	assert_line --index 0 "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: assignments add should ask for override and apply it if user replies with yes" {
	/usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>' <<< yes"
	assert_success
	assert_line --index 0 "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should force override if the '--force-override' option is used" {
	/usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team assignments add --force-override noujz 'Mr. Noujz <noujz@mr.se>'"
	assert_success
	assert_line --index 0 "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should force override if the '-f' option is used" {
	/usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team assignments add -f noujz 'Mr. Noujz <noujz@mr.se>'"
	assert_success
	assert_line --index 0 "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: assignments add should ask for override and abort if user replies with no" {
	/usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>' <<< no"
	assert_success
	assert_line --index 0 "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] Nothing changed"
}

@test "git-team: assignments add should ask for override and abort if user replies with anything else" {
	/usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>' <<< foo"
	assert_success
	assert_line --index 0 "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] Nothing changed"
}

@test "git-team: assignments add should ask for override and abort if user just hits ENTER" {
	/usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
	run bash -c "/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>' <<< ''"
	assert_success
	assert_line --index 0 "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] Nothing changed"
}

@test "git-team: assignments add should fail to create an assigment for an invalidly formatted co-author" {
	run /usr/local/bin/git-team assignments add noujz INVALID-CO-AUTHOR
	assert_failure 255
	assert_line --index 0 "error: Not a valid coauthor: INVALID-CO-AUTHOR"
}

