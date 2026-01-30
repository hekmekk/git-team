#!/usr/bin/env bats

REPO_PATH=/tmp/repo/status-tests
REPO_CHECKSUM=$(echo -n $USER:$REPO_PATH | md5sum | awk '{ print $1 }')

setup() {
	bats_load_library bats-support
	bats_load_library bats-assert

	/usr/local/bin/git-team config activation-scope repo-local

	mkdir -p $REPO_PATH
	cd $REPO_PATH

	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz
}

teardown() {
	/usr/local/bin/git-team config activation-scope global

	cd -
	rm -rf $REPO_PATH
}

@test 'git-team: (scope: repo-local) status should properly display the disabled status' {
	run /usr/local/bin/git-team status
	assert_success
	assert_line 'git-team disabled'
}

@test 'git-team: (scope: repo-local) status should properly display the disabled status in a json format' {
	run /usr/local/bin/git-team status --json
	assert_success
	assert_line --index 0 '{"status":"disabled","coAuthors":[]}'
}

@test 'git-team: (scope: repo-local) status should properly display the enabled status' {
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

@test 'git-team: (scope: repo-local) status should properly display the enabled status in a json format' {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'

	run /usr/local/bin/git-team status --json
	assert_success
	assert_line --index 0 '{"status":"enabled","coAuthors":["A <a@x.y>","B <b@x.y>","C <c@x.y>"]}'

	/usr/local/bin/git-team disable
}

@test 'git-team: (scope: repo-local) status should fail when not inside a git repository' {
	cd /tmp

	run /usr/local/bin/git-team status
	assert_failure 1
	assert_line 'error: failed to get status with activation-scope=repo-local: not inside a git repository'
}


