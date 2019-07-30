#!/usr/bin/env bats

setup() {
	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add b 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'

	/usr/local/bin/git-team enable a b c
}

@test "status enabled" {
  run /usr/local/bin/git-team status
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "git-team enabled." ]
  [ "${lines[1]}" = "Co-authors:" ]
  [ "${lines[2]}" = "-----------" ]
  [ "${lines[3]}" = "A <a@x.y>" ]
  [ "${lines[4]}" = "B <b@x.y>" ]
  [ "${lines[5]}" = "C <c@x.y>" ]
}

teardown() {
	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm b
	/usr/local/bin/git-team rm c

	/usr/local/bin/git-team disable
}


