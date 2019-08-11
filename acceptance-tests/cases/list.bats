#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

setup() {
	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add b 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'
}

@test "git-team: list should show all alias -> coauthor assignments" {
  run /usr/local/bin/git-team list
  assert_success
  assert_line --index 0 'Aliases:'
  assert_line --index 1 '--------' ]
  assert_line --index 2 "'a' -> 'A <a@x.y>'"
  assert_line --index 3 "'b' -> 'B <b@x.y>'"
  assert_line --index 4 "'c' -> 'C <c@x.y>'"
}

@test "git-team: ls should show all alias -> coauthor assignments" {
  run /usr/local/bin/git-team ls
  assert_success
  assert_line --index 0 'Aliases:'
  assert_line --index 1 '--------' ]
  assert_line --index 2 "'a' -> 'A <a@x.y>'"
  assert_line --index 3 "'b' -> 'B <b@x.y>'"
  assert_line --index 4 "'c' -> 'C <c@x.y>'"
}

teardown() {
	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm b
	/usr/local/bin/git-team rm c
}


