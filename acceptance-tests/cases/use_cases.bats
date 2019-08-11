#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

LOCAL_CONFIG_PATH=/root/.config/git-team
HOOKS_PATH=/usr/local/share/.config/git-team/hooks
REPO_PATH=/tmp/repo

setup() {
	mkdir -p $LOCAL_CONFIG_PATH
	mkdir -p $HOOKS_PATH
	cp /usr/local/bin/prepare-commit-msg $HOOKS_PATH/prepare-commit-msg

	mkdir -p $REPO_PATH
}

teardown() {
	rm -rf $LOCAL_CONFIG_PATH
	rm -rf $HOOKS_PATH
	rm -rf $REPO_PATH
}

@test "use case: when git-team is enabled then 'git commit -m' should have the respective co-authors injected" {
	/usr/local/bin/git-team enable 'B <b@x.y>' 'A <a@x.y>' 'C <c@x.y>'

	cd $REPO_PATH
	touch THE_FILE

	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz

	git add -A
	git commit -m "test"

	run git show --name-only
	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 'Author: git-team-acceptance-test <foo@bar.baz>'
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+test'
	refute_line --index 4 --regexp '\w+'
	assert_line --index 5 --regexp '\s+Co-authored-by: A <a@x.y>'
	assert_line --index 6 --regexp '\s+Co-authored-by: B <b@x.y>'
	assert_line --index 7 --regexp '\s+Co-authored-by: C <c@x.y>'
	assert_line --index 8 'THE_FILE'

	cd -
	/usr/local/bin/git-team disable
}

@test "use case: when git-team is disabled then 'git commit -m' should not have any co-authors injected" {
	/usr/local/bin/git-team disable

	cd $REPO_PATH
	touch THE_FILE

	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz

	git add -A
	git commit -m "test"

	run git show --name-only
	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 'Author: git-team-acceptance-test <foo@bar.baz>'
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+test'
	assert_line --index 4 'THE_FILE'
	refute_output --partial 'Co-authored-by:'

	cd -
}
