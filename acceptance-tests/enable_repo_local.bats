#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

REPO_PATH=/tmp/repo/enable-tests
REPO_CHECKSUM=$(echo -n $USER:$REPO_PATH | md5sum | awk '{ print $1 }')

setup() {
	cp /usr/local/bin/prepare-commit-msg /usr/local/etc/git-team/hooks/prepare-commit-msg

	/usr/local/bin/git-team config activation-scope repo-local

	mkdir -p $REPO_PATH
	cd $REPO_PATH

	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz

	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add b 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'
}

teardown() {
	/usr/local/bin/git-team disable

	/usr/local/bin/git-team config activation-scope global

	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm b
	/usr/local/bin/git-team rm c

	cd -
	rm -rf $REPO_PATH
}

@test "git-team: (scope: repo-local) enable should persist the current status" {
	/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>'

	run bash -c "git config --local --get-regexp team.state | sort"
	assert_success
	assert_line --index 0 'team.state.active-coauthors A <a@x.y>'
	assert_line --index 1 'team.state.active-coauthors Ad-hoc <adhoc@tmp.se>'
	assert_line --index 2 'team.state.active-coauthors B <b@x.y>'
	assert_line --index 3 'team.state.active-coauthors C <c@x.y>'
	assert_line --index 4 'team.state.status enabled'
}

@test "git-team: (scope: repo-local) enable should set the prepare-commit-msg hook" {
	run bash -c "/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>' &>/dev/null && git config --local core.hooksPath"
	assert_success
	assert_line '/usr/local/etc/git-team/hooks'
}

@test "git-team: (scope: repo-local) enable should set the commit template" {
	run bash -c "/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>' &>/dev/null && git config --local commit.template"
	assert_success
	assert_line "/root/.config/git-team/commit-templates/repo-local/$REPO_CHECKSUM/COMMIT_TEMPLATE"
}

@test "git-team: (scope: repo-local) enable should provision the commit template" {
	run bash -c "/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>' &>/dev/null && cat /root/.config/git-team/commit-templates/repo-local/$REPO_CHECKSUM/COMMIT_TEMPLATE"
	assert_success
	assert_line --index 0 'Co-authored-by: A <a@x.y>'
	assert_line --index 1 'Co-authored-by: Ad-hoc <adhoc@tmp.se>'
	assert_line --index 2 'Co-authored-by: B <b@x.y>'
	assert_line --index 3 'Co-authored-by: C <c@x.y>'
}

@test "git-team: (scope: repo-local) enable shorthand should display the enabled co-authors in alphabetical order" {
	run /usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>'
	assert_success
	assert_line --index 0 'git-team enabled'
	assert_line --index 1 'co-authors'
	assert_line --index 2 '─ A <a@x.y>'
	assert_line --index 3 '─ Ad-hoc <adhoc@tmp.se>'
	assert_line --index 4 '─ B <b@x.y>'
	assert_line --index 5 '─ C <c@x.y>'
}

@test "git-team: (scope: repo-local) should display the enabled co-authors in alphabetical order" {
	run /usr/local/bin/git-team enable b a c 'Ad-hoc <adhoc@tmp.se>'
	assert_success
	assert_line --index 0 'git-team enabled'
	assert_line --index 1 'co-authors'
	assert_line --index 2 '─ A <a@x.y>'
	assert_line --index 3 '─ Ad-hoc <adhoc@tmp.se>'
	assert_line --index 4 '─ B <b@x.y>'
	assert_line --index 5 '─ C <c@x.y>'
}

@test "git-team: (scope: repo-local) issuing enable should be idempotent" {
	/usr/local/bin/git-team enable b a c 'Ad-hoc <adhoc@tmp.se>'
	/usr/local/bin/git-team enable b a c 'Ad-hoc <adhoc@tmp.se>'
	run /usr/local/bin/git-team enable b a c 'Ad-hoc <adhoc@tmp.se>'
	assert_success
	assert_line --index 0 'git-team enabled'
	assert_line --index 1 'co-authors'
	assert_line --index 2 '─ A <a@x.y>'
	assert_line --index 3 '─ Ad-hoc <adhoc@tmp.se>'
	assert_line --index 4 '─ B <b@x.y>'
	assert_line --index 5 '─ C <c@x.y>'
	assert_line --index 6 ''
	assert_line --index 7 ''
	assert_line --index 8 ''
}

@test "git-team: (scope: repo-local) enable should ignore duplicates" {
	run /usr/local/bin/git-team enable a a 'A <a@x.y>'
	assert_success
	assert_line --index 0 'git-team enabled'
	assert_line --index 1 'co-authors'
	assert_line --index 2 '─ A <a@x.y>'
	assert_line --index 3 ''
}

@test "git-team: (scope: repo-local) enable should fail when trying to enable with a non-existing alias" {
	run /usr/local/bin/git-team enable non-existing-alias
	assert_failure 255
	assert_line 'error: Failed to resolve alias team.alias.non-existing-alias'
}

@test "git-team: (scope: repo-local) enable should fail when trying to enable while not inside a git repository" {
	cd /tmp
	run /usr/local/bin/git-team enable a
	assert_failure 255
	assert_line 'error: Failed to enable with activation-scope=repo-local: not inside a git repository'
}

