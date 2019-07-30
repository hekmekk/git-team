#!/usr/bin/env bats

setup() {
	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add b 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'
}

@test "list" {
  run /usr/local/bin/git-team list
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "Aliases:" ]
  [ "${lines[1]}" = "--------" ]
  [ "${lines[2]}" = "'a' -> 'A <a@x.y>'" ]
  [ "${lines[3]}" = "'b' -> 'B <b@x.y>'" ]
  [ "${lines[4]}" = "'c' -> 'C <c@x.y>'" ]
}

@test "ls" {
  run /usr/local/bin/git-team ls
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "Aliases:" ]
  [ "${lines[1]}" = "--------" ]
  [ "${lines[2]}" = "'a' -> 'A <a@x.y>'" ]
  [ "${lines[3]}" = "'b' -> 'B <b@x.y>'" ]
  [ "${lines[4]}" = "'c' -> 'C <c@x.y>'" ]
}

teardown() {
	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm b
	/usr/local/bin/git-team rm c
}


