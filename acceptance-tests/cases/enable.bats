#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

setup() {
	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add b 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'
}

@test "git-team: enable should persist the current status to status file" {
	run bash -c "/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>' &>/dev/null && cat /root/.config/git-team/status.toml"
	assert_success
	assert_line --index 0 'status = "enabled"'
	assert_line --index 1 'co-authors = ["A <a@x.y>", "Ad-hoc <adhoc@tmp.se>", "B <b@x.y>", "C <c@x.y>"]'
}

@test "git-team: enable should enable the prepare-commit-msg hook" {
	run bash -c "/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>' &>/dev/null && git config --global core.hooksPath"
	assert_success
	assert_line '/usr/local/share/.config/git-team/hooks'
}

@test "git-team: enable should set a commit template" {
	run bash -c "/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>' &>/dev/null && git config --global commit.template"
	assert_success
	assert_line '/root/.config/git-team/COMMIT_TEMPLATE'
}

@test "git-team: enable should provision a commit template" {
	run bash -c "/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>' &>/dev/null && cat /root/.config/git-team/COMMIT_TEMPLATE"
	assert_success
	assert_line --index 0 'Co-authored-by: A <a@x.y>'
	assert_line --index 1 'Co-authored-by: Ad-hoc <adhoc@tmp.se>'
	assert_line --index 2 'Co-authored-by: B <b@x.y>'
	assert_line --index 3 'Co-authored-by: C <c@x.y>'
}

@test "git-team: enable shorthand should display the enabled co-authors in alphabetical order" {
	run /usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>'
	assert_success
	assert_line --index 0 'git-team enabled.'
	assert_line --index 1 'Co-authors:'
	assert_line --index 2 '-----------'
	assert_line --index 3 'A <a@x.y>'
	assert_line --index 4 'Ad-hoc <adhoc@tmp.se>'
	assert_line --index 5 'B <b@x.y>'
	assert_line --index 6 'C <c@x.y>'
}

@test "git-team: enable should display the enabled co-authors in alphabetical order" {
	run /usr/local/bin/git-team enable b a c 'Ad-hoc <adhoc@tmp.se>'
	assert_success
	assert_line --index 0 'git-team enabled.'
	assert_line --index 1 'Co-authors:'
	assert_line --index 2 '-----------'
	assert_line --index 3 'A <a@x.y>'
	assert_line --index 4 'Ad-hoc <adhoc@tmp.se>'
	assert_line --index 5 'B <b@x.y>'
	assert_line --index 6 'C <c@x.y>'
}

@test "git-team: enable should ignore duplicates" {
	run /usr/local/bin/git-team enable a a 'A <a@x.y>'
	assert_success
	assert_line --index 0 'git-team enabled.'
	assert_line --index 1 'Co-authors:'
	assert_line --index 2 '-----------'
	assert_line --index 3 'A <a@x.y>'
	assert_line --index 4 ''
}

@test "git-team: enable should fail when trying to enable with a non-existing alias" {
	run /usr/local/bin/git-team enable non-existing-alias
	assert_failure 255
	assert_line 'error: Failed to resolve alias team.alias.non-existing-alias'
}

teardown() {
	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm b
	/usr/local/bin/git-team rm c

	/usr/local/bin/git-team disable
}

