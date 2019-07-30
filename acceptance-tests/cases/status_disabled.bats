#!/usr/bin/env bats

@test "status disabled" {
  run /usr/local/bin/git-team status
  [ "$status" -eq 0 ]
  [ "$output" = "git-team disabled." ]
}
