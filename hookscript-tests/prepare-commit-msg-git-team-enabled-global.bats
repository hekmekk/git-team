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

REPO_PATH=/tmp/repo/prepare-commit-msg-enabled-global

setup() {
  bats_load_library bats-support
  bats_load_library bats-assert

  touch /tmp/COMMIT_MSG
  mkdir -p $REPO_PATH
  cd $REPO_PATH
  git init
  git config user.name git-team-acceptance-test
  git config user.email foo@bar.baz
  /usr/local/bin/git-team config activation-scope global
  /usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'
}

teardown() {
  /usr/local/bin/git-team disable
  cd -
  rm -rf $REPO_PATH
  rm /tmp/COMMIT_MSG
}

@test "prepare-commit-msg: git-team enabled: (scope: global) - message" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG message && cat /tmp/COMMIT_MSG"
  assert_success
  assert_line --index 0 'Co-authored-by: A <a@x.y>'
  assert_line --index 1 'Co-authored-by: B <b@x.y>'
  assert_line --index 2 'Co-authored-by: C <c@x.y>'
}

@test "prepare-commit-msg: git-team enabled: (scope: global) - none" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG && cat /tmp/COMMIT_MSG"
  assert_success
  refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team enabled: (scope: global) - commit" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG commit && cat /tmp/COMMIT_MSG"
  assert_success
  assert_line --index 0 'Co-authored-by: A <a@x.y>'
  assert_line --index 1 'Co-authored-by: B <b@x.y>'
  assert_line --index 2 'Co-authored-by: C <c@x.y>'
}

@test "prepare-commit-msg: git-team enabled: (scope: global) - template" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG template && cat /tmp/COMMIT_MSG"
  assert_success
  refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team enabled: (scope: global) - merge" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG merge && cat /tmp/COMMIT_MSG"
  assert_success
  assert_line --index 0 'Co-authored-by: A <a@x.y>'
  assert_line --index 1 'Co-authored-by: B <b@x.y>'
  assert_line --index 2 'Co-authored-by: C <c@x.y>'
}

@test "prepare-commit-msg: git-team enabled: (scope: global) - squash" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG squash && cat /tmp/COMMIT_MSG"
  assert_success
  assert_line --index 0 'Co-authored-by: A <a@x.y>'
  assert_line --index 1 'Co-authored-by: B <b@x.y>'
  assert_line --index 2 'Co-authored-by: C <c@x.y>'
}
