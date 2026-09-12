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

REPO_PATH=/tmp/repo/use-cases-global-local-interaction
USER_NAME=git-team-acceptance-test
USER_EMAIL=acc@git.team

setup() {
  bats_load_library bats-support
  bats_load_library bats-assert

  git config --global init.defaultBranch main
  git config --global user.name "$USER_NAME"
  git config --global user.email "$USER_EMAIL"

  mkdir -p $REPO_PATH
  cd $REPO_PATH

  git init
}

teardown() {
  /usr/local/bin/git-team config activation-scope global
  /usr/local/bin/git-team disable

  /usr/local/bin/git-team config activation-scope repo-local
  /usr/local/bin/git-team disable

  cd -
  rm -rf $REPO_PATH

  rm /home/git-team-acceptance-test/.gitconfig
}

@test "use case: enabled global should not interfere with disabled repo-local" {
  /usr/local/bin/git-team config activation-scope global
  /usr/local/bin/git-team enable 'A <a@x.y>'

  /usr/local/bin/git-team config activation-scope repo-local
  /usr/local/bin/git-team disable

  /usr/local/bin/git-team config activation-scope global

  run git config --global core.hooksPath
  assert_success
  assert_output '/home/git-team-acceptance-test/.git-team/hooks'

  run git config --global commit.template
  assert_success
  assert_output '/home/git-team-acceptance-test/.git-team/commit-templates/global/COMMIT_TEMPLATE'

  run git config core.hooksPath
  assert_success
  assert_output '/home/git-team-acceptance-test/.git-team/hooks'

  run git config commit.template
  assert_success
  assert_output '/home/git-team-acceptance-test/.git-team/commit-templates/global/COMMIT_TEMPLATE'

  touch THE_FILE_GLOBAL
  git add -A
  git commit -m "test global enabled"

  run git show --name-only
  assert_success
  assert_output --partial 'Co-authored-by:'

  /usr/local/bin/git-team config activation-scope repo-local

  touch THE_FILE_REPO_LOCAL
  git add -A
  git commit -m "test repo-local disabled"

  run git show --name-only
  assert_success
  refute_output --partial 'Co-authored-by:'
}

@test "use case: enabled repo-local should not interfere with disabled global" {
  /usr/local/bin/git-team config activation-scope repo-local
  /usr/local/bin/git-team enable 'A <a@x.y>'

  /usr/local/bin/git-team config activation-scope global
  /usr/local/bin/git-team disable

  /usr/local/bin/git-team config activation-scope repo-local

  touch THE_FILE_GLOBAL
  git add -A
  git commit -m "test local enabled"

  run git show --name-only
  assert_success
  assert_output --partial 'Co-authored-by:'

  /usr/local/bin/git-team config activation-scope global

  touch THE_FILE_REPO_LOCAL
  git add -A
  git commit -m "test global disabled"

  run git show --name-only
  assert_success
  refute_output --partial 'Co-authored-by:'
}

@test "use case: disabled global should not interfere with enabled repo-local" {
  /usr/local/bin/git-team config activation-scope global
  /usr/local/bin/git-team disable

  /usr/local/bin/git-team config activation-scope repo-local
  /usr/local/bin/git-team enable 'A <a@x.y>'

  /usr/local/bin/git-team config activation-scope global

  touch THE_FILE_GLOBAL
  git add -A
  git commit -m "test global disabled"

  run git show --name-only
  assert_success
  refute_output --partial 'Co-authored-by:'

  /usr/local/bin/git-team config activation-scope repo-local

  touch THE_FILE_REPO_LOCAL
  git add -A
  git commit -m "test repo-local enabled"

  run git show --name-only
  assert_success
  assert_output --partial 'Co-authored-by:'
}

@test "use case: disabled repo-local should not interfere with enabled global" {
  /usr/local/bin/git-team config activation-scope repo-local
  /usr/local/bin/git-team disable

  /usr/local/bin/git-team config activation-scope global
  /usr/local/bin/git-team enable 'A <a@x.y>'

  /usr/local/bin/git-team config activation-scope repo-local

  touch THE_FILE_REPO_LOCAL
  git add -A
  git commit -m "test repo-local disabled"

  run git show --name-only
  assert_success
  refute_output --partial 'Co-authored-by:'

  /usr/local/bin/git-team config activation-scope global

  touch THE_FILE_GLOBAL
  git add -A
  git commit -m "test global enabled"

  run git show --name-only
  assert_success
  assert_output --partial 'Co-authored-by:'
}
