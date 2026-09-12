#!/usr/bin/env bats

# Copyright (C) 2026 Rea Sand

# This file is part of git-team.

# This program is free software: you can redistribute it and/or
# modify it under the terms of the GNU General Public License
# as published by the Free Software Foundation, either version 3
# of the License, or (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program. If not, see <https://www.gnu.org/licenses/>.

setup() {
  bats_load_library bats-support
  bats_load_library bats-assert
}

teardown() {
  bash -c "/usr/local/bin/git-team assignments rm noujz || true"

  bash -c "/usr/local/bin/git-team assignments rm green || true"

  bash -c "/usr/local/bin/git-team assignments rm a || true"
  bash -c "/usr/local/bin/git-team assignments rm b || true"
  bash -c "/usr/local/bin/git-team assignments rm c || true"
}

@test "git-team: assignments add should add an assignment to git config" {
  run bash -c "/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>' &>/dev/null && git config --global team.alias.noujz"
  assert_success
  assert_line 'Mr. Noujz <noujz@mr.se>'
}

@test "git-team: assignments add should create a new assigment" {
  run /usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>'
  assert_success
  assert_line --index 0 "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: assignments add should create a new assigment when receiving input from stdin" {
  run bash -c "echo noujz 'Mr. Noujz <noujz@mr.se>' | /usr/local/bin/git-team assignments add"
  assert_success
  assert_line "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should create an assigment when receiving multiple lines of input from stdin" {
  run bash -c "for alias in a b c; do echo \$alias 'Mrs. Noujz <noujz@mrs.se>'; done | /usr/local/bin/git-team assignments add"
  assert_success
  assert_line --index 0 "Assignment added: 'a' →  'Mrs. Noujz <noujz@mrs.se>'"
  assert_line --index 1 "Assignment added: 'b' →  'Mrs. Noujz <noujz@mrs.se>'"
  assert_line --index 2 "Assignment added: 'c' →  'Mrs. Noujz <noujz@mrs.se>'"
}

@test "git-team: assignments add should ask for override and apply it if user replies with yes" {
  /usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
  run bash -c "/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>' <<< yes"
  assert_success
  assert_line --index 0 "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should force override if the '--force-override' option is used" {
  /usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
  run bash -c "/usr/local/bin/git-team assignments add --force-override noujz 'Mr. Noujz <noujz@mr.se>'"
  assert_success
  assert_line --index 0 "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should force override if the '-f' option is used" {
  /usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
  run bash -c "/usr/local/bin/git-team assignments add -f noujz 'Mr. Noujz <noujz@mr.se>'"
  assert_success
  assert_line --index 0 "Assignment added: 'noujz' →  'Mr. Noujz <noujz@mr.se>'"
}

@test "git-team: add should keep the existing assignment if the '--keep-existing' option is used" {
  /usr/local/bin/git-team assignments add green 'Green <green@git.team>'
  run bash -c "/usr/local/bin/git-team assignments add --keep-existing green 'Red <red@git.team>'"
  assert_success
  refute_output --regexp '\w+'
}

@test "git-team: add should keep the existing assignment if the '-k' option is used" {
  /usr/local/bin/git-team assignments add green 'Green <green@git.team>'
  run bash -c "/usr/local/bin/git-team assignments add -k green 'Red <red@git.team>'"
  assert_success
  refute_output --regexp '\w+'
}

@test "git-team: assignments add should ask for override and abort if user replies with no" {
  /usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
  run bash -c "/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>' <<< no"
  assert_success
  assert_line --index 0 "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] "
}

@test "git-team: assignments add should ask for override and abort if user replies with anything else" {
  /usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
  run bash -c "/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>' <<< foo"
  assert_success
  assert_line --index 0 "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] "
}

@test "git-team: assignments add should ask for override and abort if user just hits ENTER" {
  /usr/local/bin/git-team assignments add noujz 'Mr. Green <green@mr.se>'
  run bash -c "/usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>' <<< ''"
  assert_success
  assert_line --index 0 "Assignment 'noujz' →  'Mr. Green <green@mr.se>' exists already. Override with 'Mr. Noujz <noujz@mr.se>'? [N/y] "
}

@test "git-team: assignments add should fail to create an assigment for an invalidly formatted co-author" {
  run /usr/local/bin/git-team assignments add noujz INVALID-CO-AUTHOR
  assert_failure 1
  assert_line --index 0 "error: not a valid coauthor: INVALID-CO-AUTHOR"
}
