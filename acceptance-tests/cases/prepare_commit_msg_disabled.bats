#!/usr/bin/env bats

setup() {
	/usr/local/bin/git-team disable
}

@test "commit -m" {
  run /usr/local/bin/prepare-commit-msg foo message
  [ "$status" -eq 0 ]
}

