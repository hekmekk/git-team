#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

@test "git-team: ls should show 'No assignments'" {
	run /usr/local/bin/git-team ls
	assert_success
	assert_line 'No assignments'
}

@test "git-team: list should show all alias -> coauthor assignments" {
	/usr/local/bin/git-team assignments add a 'A <a@x.y>'
	/usr/local/bin/git-team assignments add bb 'B <b@x.y>'
	/usr/local/bin/git-team assignments add c 'C <c@x.y>'

	run /usr/local/bin/git-team list
	assert_success
	assert_line --index 0 'assignments'
	assert_line --index 1 '─ a  →  A <a@x.y>'
	assert_line --index 2 '─ bb →  B <b@x.y>'
	assert_line --index 3 '─ c  →  C <c@x.y>'

	/usr/local/bin/git-team assignments rm a
	/usr/local/bin/git-team assignments rm bb
	/usr/local/bin/git-team assignments rm c
}

@test "git-team: ls should show all alias -> coauthor assignments" {
	/usr/local/bin/git-team assignments add a 'A <a@x.y>'
	/usr/local/bin/git-team assignments add bb 'B <b@x.y>'
	/usr/local/bin/git-team assignments add c 'C <c@x.y>'

	run /usr/local/bin/git-team ls
	assert_success
	assert_line --index 0 'assignments'
	assert_line --index 1 '─ a  →  A <a@x.y>'
	assert_line --index 2 '─ bb →  B <b@x.y>'
	assert_line --index 3 '─ c  →  C <c@x.y>'

	/usr/local/bin/git-team assignments rm a
	/usr/local/bin/git-team assignments rm bb
	/usr/local/bin/git-team assignments rm c
}

