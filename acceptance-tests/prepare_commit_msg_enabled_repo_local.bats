#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

REPO_PATH=/tmp/repo/prepare-commit-msg-enabled-repo-local
REPO_CHECKSUM=$(echo -n $USER:$REPO_PATH | md5sum | awk '{ print $1 }')

setup() {
	touch /tmp/COMMIT_MSG
	mkdir -p $REPO_PATH
	cd $REPO_PATH
	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz
	/usr/local/bin/git-team config activation-scope repo-local
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'
}

teardown() {
	/usr/local/bin/git-team disable
	/usr/local/bin/git-team config activation-scope global
	cd -
	rm -rf $REPO_PATH
	rm /tmp/COMMIT_MSG
}

@test "prepare-commit-msg: git-team enabled: (scope: repo-local) - message" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG message && cat /tmp/COMMIT_MSG"
	assert_success
	assert_line --index 0 'Co-authored-by: A <a@x.y>'
	assert_line --index 1 'Co-authored-by: B <b@x.y>'
	assert_line --index 2 'Co-authored-by: C <c@x.y>'
}

@test "prepare-commit-msg: git-team enabled: (scope: repo-local) - none" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team enabled: (scope: repo-local) - commit" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG commit && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team enabled: (scope: repo-local) - template" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG template && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team enabled: (scope: repo-local) - merge" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG merge && cat /tmp/COMMIT_MSG"
	assert_success
	assert_line --index 0 'Co-authored-by: A <a@x.y>'
	assert_line --index 1 'Co-authored-by: B <b@x.y>'
	assert_line --index 2 'Co-authored-by: C <c@x.y>'
}

@test "prepare-commit-msg: git-team enabled: (scope: repo-local) - squash" {
	run bash -c "/usr/local/bin/prepare-commit-msg /tmp/COMMIT_MSG squash && cat /tmp/COMMIT_MSG"
	assert_success
	assert_line --index 0 'Co-authored-by: A <a@x.y>'
	assert_line --index 1 'Co-authored-by: B <b@x.y>'
	assert_line --index 2 'Co-authored-by: C <c@x.y>'
}

