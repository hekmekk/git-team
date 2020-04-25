#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

REPO_PATH=/tmp/repo/disable-tests
REPO_CHECKSUM=$(echo -n $USER:$REPO_PATH | md5sum | awk '{ print $1 }')

setup() {
	mkdir -p $REPO_PATH
}

teardown() {
	rm -rf $REPO_PATH
}

@test "git-team: disable should disable a previously enabled git-team with activation-scope global" {
	mkdir -p /root/.config/git-team/

	/usr/local/bin/git-team config activation-scope global
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'

	run /usr/local/bin/git-team disable
	assert_success
	assert_line "git-team disabled"

	rm -rf /root/.config/git-team/
}

@test "git-team: disable should disable a previously enabled git-team with activation-scope repo-local" {
	mkdir -p /root/.config/git-team/

	/usr/local/bin/git-team config activation-scope repo-local

	cd $REPO_PATH
	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz

	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'

	run /usr/local/bin/git-team disable
	assert_success
	assert_line "git-team disabled"

	rm -rf /root/.config/git-team/
	cd -
}

@test "git-team: disable should persist the current status to gitconfig with activation-scope global" {
	mkdir -p /root/.config/git-team/

	/usr/local/bin/git-team config activation-scope global
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run git config --global --get-regexp team.state
	assert_success
	assert_output 'team.state.status disabled'

	rm -rf /root/.config/git-team/
}

@test "git-team: disable should persist the current status to gitconfig with activation-scope repo-local" {
	mkdir -p /root/.config/git-team/

	/usr/local/bin/git-team config activation-scope repo-local

	cd $REPO_PATH
	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz

	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run git config --local --get-regexp team.state
	assert_success
	assert_output 'team.state.status disabled'

	rm -rf /root/.config/git-team/
	cd -
}

@test "git-team: disable should disable the prepare-commit-msg hook with activation-scope global" {
	mkdir -p /root/.config/git-team/

	/usr/local/bin/git-team config activation-scope global
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "git config --global core.hooksPath"
	assert_failure 1
	refute_line --regexp '\w+'

	rm -rf /root/.config/git-team/
}

@test "git-team: disable should disable the prepare-commit-msg hook with activation-scope repo-local" {
	mkdir -p /root/.config/git-team/

	/usr/local/bin/git-team config activation-scope repo-local

	cd $REPO_PATH
	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz

	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "git config --local core.hooksPath"
	assert_failure 1
	refute_line --regexp '\w+'

	rm -rf /root/.config/git-team/
	cd -
}

@test "git-team: disable should unset the commit template with activation-scope global" {
	mkdir -p /root/.config/git-team/

	/usr/local/bin/git-team config activation-scope global
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "git config --global commit.template"
	assert_failure 1
	refute_line --regexp '\w+'

	rm -rf /root/.config/git-team/
}

@test "git-team: disable should unset the commit template with activation-scope repo-local" {
	mkdir -p /root/.config/git-team/

	/usr/local/bin/git-team config activation-scope repo-local

	cd $REPO_PATH
	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz

	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "git config --local commit.template"
	assert_failure 1
	refute_line --regexp '\w+'

	rm -rf /root/.config/git-team/
	cd -
}

@test "git-team: disable should remove the according COMMIT_TEMPLATE for activation-scope global" {
	mkdir -p /root/.config/git-team/

	/usr/local/bin/git-team config activation-scope global
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "ls -la /root/.config/git-team/commit-templates/global/COMMIT_TEMPLATE"
	assert_failure 2
	assert_line "ls: cannot access '/root/.config/git-team/commit-templates/global/COMMIT_TEMPLATE': No such file or directory"

	rm -rf /root/.config/git-team/
}

@test "git-team: disable should remove the according commit-template directory for activation-scope repo-local" {
	mkdir -p /root/.config/git-team/

	/usr/local/bin/git-team config activation-scope global

	cd $REPO_PATH
	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz

	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "ls -la /root/.config/git-team/commit-templates/repo-local/$REPO_CHECKSUM"
	assert_failure 2
	assert_line "ls: cannot access '/root/.config/git-team/commit-templates/repo-local/$REPO_CHECKSUM': No such file or directory"

	rm -rf /root/.config/git-team/
	cd -
}

@test "git-team: disable should treat a previously disabled git-team idempotently" {
	mkdir -p /root/.config/git-team/

	run /usr/local/bin/git-team disable
	assert_success
	assert_line "git-team disabled"

	rm -rf /root/.config/git-team/
}

@test "git-team: disable should fail when trying to disable with activation-scope repo-local when not in a git repository directory" {
	mkdir -p /root/.config/git-team/
	/usr/local/bin/git-team config activation-scope repo-local

	run /usr/local/bin/git-team disable
	assert_failure 255
	assert_line 'error: Failed to disable with scope=repo-local: not inside a git repository'

	/usr/local/bin/git-team config activation-scope global
	rm -rf /root/.config/git-team/
}

