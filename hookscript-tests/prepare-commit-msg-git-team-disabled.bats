#!/usr/bin/env bats

REPO_PATH=/tmp/repo/prepare-commit-msg-disabled

setup() {
  bats_load_library bats-support
  bats_load_library bats-assert

  touch /tmp/COMMIT_MSG
  mkdir -p $REPO_PATH
  cd $REPO_PATH
  git init
  git config user.name git-team-acceptance-test
  git config user.email foo@bar.baz
  /usr/local/bin/git-team enable 'A <a@x.y>' # enable once to install the scripts
  /usr/local/bin/git-team disable
}

teardown() {
  cd -
  rm -rf $REPO_PATH
  rm /tmp/COMMIT_MSG
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - message" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG message && cat /tmp/COMMIT_MSG"
  assert_success
  refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - none" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG && cat /tmp/COMMIT_MSG"
  assert_success
  refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - commit" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG commit && cat /tmp/COMMIT_MSG"
  assert_success
  refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - template" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG template && cat /tmp/COMMIT_MSG"
  assert_success
  refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - merge" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG merge && cat /tmp/COMMIT_MSG"
  assert_success
  refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - squash" {
  run bash -c "~/.git-team/hooks/prepare-commit-msg-git-team.sh /tmp/COMMIT_MSG squash && cat /tmp/COMMIT_MSG"
  assert_success
  refute_output --regexp '\w+'
}
