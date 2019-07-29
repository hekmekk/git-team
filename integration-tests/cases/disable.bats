#!/usr/bin/env bats

setup() {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
}

@test "disable" {
  run /usr/local/bin/git-team disable
  [ "$status" -eq 0 ]
  [ "$output" = "git-team disabled." ]
}

