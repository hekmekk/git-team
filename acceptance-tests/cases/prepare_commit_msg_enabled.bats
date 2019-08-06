#!/usr/bin/env bats

setup() {
	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add b 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'

	/usr/local/bin/git-team enable a b c
	touch /tmp/foo
}

@test "prepare-commit-msg: git-team enabled - commit -m" {
  run /usr/local/bin/prepare-commit-msg /tmp/foo message && cat /tmp/foo
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "" ]
  [ "${lines[1]}" = "" ]
  [ "${lines[2]}" = "Co-authored-by: A <a@x.y>" ]
  [ "${lines[3]}" = "Co-authored-by: B <b@x.y>" ]
  [ "${lines[4]}" = "Co-authored-by: C <c@x.y>" ]
}

teardown() {
	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm b
	/usr/local/bin/git-team rm c

	/usr/local/bin/git-team disable
	rm /tmp/foo
}


