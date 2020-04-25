#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

REPO_PATH=/tmp/repo/status-tests
REPO_CHECKSUM=$(echo -n $USER:$REPO_PATH | md5sum | awk '{ print $1 }')

setup() {
	mkdir -p $REPO_PATH
}

teardown() {
	/usr/local/bin/git-team config activation-scope global
	rm -rf $REPO_PATH
}

@test 'git-team: status should properly display a disabled status for global activation-scope' {
	/usr/local/bin/git-team config activation-scope global

	run /usr/local/bin/git-team status
	assert_success
	assert_line 'git-team disabled'
}

@test 'git-team: status should properly display a disabled status for repo-local activation-scope' {
	/usr/local/bin/git-team config activation-scope repo-local
	cd $REPO_PATH
	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz

	run /usr/local/bin/git-team status
	assert_success
	assert_line 'git-team disabled'

	/usr/local/bin/git-team config activation-scope global
	cd -
}

@test 'git-team: status should properly disaplay the enabled status for global activation-scope' {
	/usr/local/bin/git-team config activation-scope global
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'

	run /usr/local/bin/git-team status
	assert_success
	assert_line --index 0 'git-team enabled'
	assert_line --index 1 'co-authors'
	assert_line --index 2 '─ A <a@x.y>'
	assert_line --index 3 '─ B <b@x.y>'
	assert_line --index 4 '─ C <c@x.y>'

	/usr/local/bin/git-team disable
}

@test 'git-team: status should properly disaplay the enabled status for repo-local activation-scope' {
	/usr/local/bin/git-team config activation-scope repo-local
	cd $REPO_PATH
	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'

	run /usr/local/bin/git-team status
	assert_success
	assert_line --index 0 'git-team enabled'
	assert_line --index 1 'co-authors'
	assert_line --index 2 '─ A <a@x.y>'
	assert_line --index 3 '─ B <b@x.y>'
	assert_line --index 4 '─ C <c@x.y>'

	/usr/local/bin/git-team disable
	/usr/local/bin/git-team config activation-scope global
	cd -
}


@test "git-team: status should fail for activation-scope repo-local when not in a git repository directory" {
	/usr/local/bin/git-team config activation-scope repo-local

	run /usr/local/bin/git-team status
	assert_failure 255
	assert_line 'error: Failed to get status with activation-scope=repo-local: not inside a git repository'

	/usr/local/bin/git-team config activation-scope global
}

