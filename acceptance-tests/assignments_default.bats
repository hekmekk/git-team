#!/usr/bin/env bats

setup() {
	bats_load_library bats-support
	bats_load_library bats-assert
}

@test "git-team: assignments (default) should show 'no assignments'" {
	run /usr/local/bin/git-team assignments
	assert_success
	assert_line --index 0 'No assignments'
}

@test "git-team: assignments (default) should show all alias -> coauthor assignments" {
	/usr/local/bin/git-team assignments add a 'A <a@x.y>'
	/usr/local/bin/git-team assignments add bb 'B <b@x.y>'
	/usr/local/bin/git-team assignments add c 'C <c@x.y>'

	run /usr/local/bin/git-team assignments
	assert_success
	assert_line --index 0 'Assignments'
	assert_line --index 1 '─ a  →  A <a@x.y>'
	assert_line --index 2 '─ bb →  B <b@x.y>'
	assert_line --index 3 '─ c  →  C <c@x.y>'

	/usr/local/bin/git-team assignments rm a
	/usr/local/bin/git-team assignments rm bb
	/usr/local/bin/git-team assignments rm c
}

@test "git-team: assignments ls should show all alias -> coauthor assignments" {
	/usr/local/bin/git-team assignments add a 'A <a@x.y>'
	/usr/local/bin/git-team assignments add bb 'B <b@x.y>'
	/usr/local/bin/git-team assignments add c 'C <c@x.y>'

	run /usr/local/bin/git-team assignments ls
	assert_success
	assert_line --index 0 'Assignments'
	assert_line --index 1 '─ a  →  A <a@x.y>'
	assert_line --index 2 '─ bb →  B <b@x.y>'
	assert_line --index 3 '─ c  →  C <c@x.y>'

	/usr/local/bin/git-team assignments rm a
	/usr/local/bin/git-team assignments rm bb
	/usr/local/bin/git-team assignments rm c
}

