#!/usr/bin/env bats

@test "add success" {
  run /usr/local/bin/git-team add noujz 'Mr. Noujz <noujz@mr.se>'
  [ "$status" -eq 0 ]
  [ "$output" = "Alias 'noujz' -> 'Mr. Noujz <noujz@mr.se>' has been added." ]
}

@test "add failure" {
  run /usr/local/bin/git-team add noujz INVALID-CO-AUTHOR
  [ "$status" -eq 255 ]
  [ "$output" = "error: Not a valid coauthor: INVALID-CO-AUTHOR" ]
}

