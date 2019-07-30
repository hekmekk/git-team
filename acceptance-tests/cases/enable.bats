#!/usr/bin/env bats

setup() {
	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add b 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'
}

@test "enable success" {
  run /usr/local/bin/git-team enable a b c 'Ad-hoc <adhoc@tmp.se>'
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "git-team enabled." ]
  [ "${lines[1]}" = "Co-authors:" ]
  [ "${lines[2]}" = "-----------" ]
  [ "${lines[3]}" = "Ad-hoc <adhoc@tmp.se>" ]
  [ "${lines[4]}" = "A <a@x.y>" ]
  [ "${lines[5]}" = "B <b@x.y>" ]
  [ "${lines[6]}" = "C <c@x.y>" ]
}

@test "enable failure" {
  run /usr/local/bin/git-team enable non-existing-alias
  [ "$status" -eq 255 ]
  [ "$output" = "error: Failed to resolve alias team.alias.non-existing-alias" ]
}

teardown() {
	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm b
	/usr/local/bin/git-team rm c
}

