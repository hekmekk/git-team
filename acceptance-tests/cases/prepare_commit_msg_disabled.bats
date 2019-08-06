#!/usr/bin/env bats

setup() {
	/usr/local/bin/git-team disable
	touch /tmp/foo
}

@test "prepare-commit-msg: git-team disabled - commit -m" {
  run /usr/local/bin/prepare-commit-msg /tmp/foo message && cat /tmp/foo
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "" ]
  [ "${lines[1]}" = "" ]
  [ "${lines[2]}" = "" ]
}

teardown() {
	rm /tmp/foo
}

