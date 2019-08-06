#!/usr/bin/env bats

setup() {
	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add b 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'

	/usr/local/bin/git-team enable a b c
}

@test "commit -m" {
  run /usr/local/bin/prepare-commit-msg foo message
  [ "$status" -eq 0 ]
}

teardown() {
	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm b
	/usr/local/bin/git-team rm c

	/usr/local/bin/git-team disable
}


