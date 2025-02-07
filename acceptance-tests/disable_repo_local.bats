#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

REPO_PATH=/tmp/repo/disable-tests
REPO_CHECKSUM=$(echo -n $USER:$REPO_PATH | md5sum | awk '{ print $1 }')
USER_NAME=git-team-acceptance-test
USER_EMAIL=acc@git.team

setup() {
	git config --global init.defaultBranch main

	mkdir -p $REPO_PATH
	cd $REPO_PATH
	git init
	git config user.name "$USER_NAME"
	git config user.email "$USER_EMAIL"

	/usr/local/bin/git-team config activation-scope repo-local
}

teardown() {
	/usr/local/bin/git-team config activation-scope global

	cd -
	rm -rf $REPO_PATH

	rm /home/git-team-acceptance-test/.gitconfig
}

@test "git-team: (scope: repo-local) disable should disable a previously enabled git-team" {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'

	run /usr/local/bin/git-team disable
	assert_success
	assert_line "git-team disabled"
}

@test "git-team: (scope: repo-local) disable should persist the current status to gitconfig" {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run git config --local --get-regexp team.state
	assert_success
	assert_output 'team.state.status disabled'
}

@test "git-team: (scope: repo-local) disable should persist a previous hooks path as the current hooks path" {
	git config --local core.hooksPath "/path/to/non-git-team-hooks"
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "git config --local core.hooksPath"
	assert_success
	assert_line --index 0 '/path/to/non-git-team-hooks'

	git config --local --unset core.hooksPath | true
}

@test "git-team: (scope: repo-local) disable should disable the prepare-commit-msg hook" {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "git config --local core.hooksPath"
	assert_failure 1
	refute_line --regexp '\w+'
}

@test "git-team: (scope: repo-local) disable should unset the commit template" {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "git config --local commit.template"
	assert_failure 1
	refute_line --regexp '\w+'
}

@test "git-team: (scope: repo-local) disable should remove the according commit-template directory" {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "ls -la /home/git-team-acceptance-test/.git-team/commit-templates/repo-local/$REPO_CHECKSUM"
	assert_failure 1
	assert_line "ls: /home/git-team-acceptance-test/.git-team/commit-templates/repo-local/$REPO_CHECKSUM: No such file or directory"
}

@test "git-team: (scope: repo-local) disable should treat a previously disabled git-team in an idempotent way" {
	run /usr/local/bin/git-team disable
	assert_success
	assert_line "git-team disabled"
}

@test "git-team: (scope: repo-local) disable should fail when trying to disable when not in a git repository directory" {
	cd /tmp
	run /usr/local/bin/git-team disable
	assert_failure 1
	assert_line 'error: failed to disable with activation-scope=repo-local: not inside a git repository'
	cd $REPO_PATH
}

