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

@test "git-team: assignments rm should remove an assigment from git config" {
  /usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>'

  run bash -c "/usr/local/bin/git-team assignments rm noujz &>/dev/null && git config --global team.alias.noujz"
  assert_failure 1
  refute_output --regexp '\w+'
}

@test "git-team: assignments rm should remove an assigment" {
  /usr/local/bin/git-team assignments add noujz 'Mr. Noujz <noujz@mr.se>'

  run /usr/local/bin/git-team assignments rm noujz
  assert_success
  assert_line --index 0 "Assignment removed: 'noujz'"
}

@test "git-team: assignments rm should fail for a non-existing alias" {
  run /usr/local/bin/git-team assignments rm noujz
  assert_failure
  assert_line --index 0 "error: no such alias: 'noujz'"
}
