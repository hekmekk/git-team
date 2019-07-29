#!/usr/bin/env bats

@test "remove" {
  run /usr/local/bin/git-team rm noujz
  [ "$status" -eq 0 ]
  [ "$output" = "Alias 'noujz' has been removed." ]
}

