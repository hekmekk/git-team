#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

@test "git-team: ls should show 'No assignments'" {
	run /usr/local/bin/git-team ls
	assert_success
	assert_line --index 0 "warn: 'git team ls' has been deprecated and is going to be removed in a future major release, use 'git team assignments' instead"
	assert_line --index 1 'No assignments'
}

@test "git-team: list should show all alias -> coauthor assignments" {
	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add bb 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'

	run /usr/local/bin/git-team list
	assert_success
	assert_line --index 0 "warn: 'git team ls' has been deprecated and is going to be removed in a future major release, use 'git team assignments' instead"
	assert_line --index 1 'assignments'
	assert_line --index 2 '─ a  →  A <a@x.y>'
	assert_line --index 3 '─ bb →  B <b@x.y>'
	assert_line --index 4 '─ c  →  C <c@x.y>'

	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm bb
	/usr/local/bin/git-team rm c
}

@test "git-team: ls should show all alias -> coauthor assignments" {
	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add bb 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'

	run /usr/local/bin/git-team ls
	assert_success
	assert_line --index 0 "warn: 'git team ls' has been deprecated and is going to be removed in a future major release, use 'git team assignments' instead"
	assert_line --index 1 'assignments'
	assert_line --index 2 '─ a  →  A <a@x.y>'
	assert_line --index 3 '─ bb →  B <b@x.y>'
	assert_line --index 4 '─ c  →  C <c@x.y>'

	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm bb
	/usr/local/bin/git-team rm c
}

