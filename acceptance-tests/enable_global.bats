#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

setup() {
	cp /usr/local/bin/prepare-commit-msg /usr/local/etc/git-team/hooks/prepare-commit-msg

	/usr/local/bin/git-team config activation-scope global

	/usr/local/bin/git-team add a 'A <a@x.y>'
	/usr/local/bin/git-team add b 'B <b@x.y>'
	/usr/local/bin/git-team add c 'C <c@x.y>'
}

teardown() {
	/usr/local/bin/git-team disable

	/usr/local/bin/git-team rm a
	/usr/local/bin/git-team rm b
	/usr/local/bin/git-team rm c
}

@test "git-team: (scope: global) enable should persist the current status to gitconfig" {
	/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>'

	run bash -c "git config --global --get-regexp team.state | sort"
	assert_success
	assert_line --index 0 'team.state.active-coauthors A <a@x.y>'
	assert_line --index 1 'team.state.active-coauthors Ad-hoc <adhoc@tmp.se>'
	assert_line --index 2 'team.state.active-coauthors B <b@x.y>'
	assert_line --index 3 'team.state.active-coauthors C <c@x.y>'
	assert_line --index 4 'team.state.status enabled'

	/usr/local/bin/git-team disable
}

@test "git-team: (scope: global) enable should enable the prepare-commit-msg hook" {
	run bash -c "/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>' &>/dev/null && git config --global core.hooksPath"
	assert_success
	assert_line '/usr/local/etc/git-team/hooks'
}

@test "git-team: (scope: global) enable should set a commit template" {
	run bash -c "/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>' &>/dev/null && git config --global commit.template"
	assert_success
	assert_line '/root/.config/git-team/commit-templates/global/COMMIT_TEMPLATE'
}

@test "git-team: (scope: global) enable should provision the commit template" {
	run bash -c "/usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>' &>/dev/null && cat /root/.config/git-team/commit-templates/global/COMMIT_TEMPLATE"
	assert_success
	assert_line --index 0 'Co-authored-by: A <a@x.y>'
	assert_line --index 1 'Co-authored-by: Ad-hoc <adhoc@tmp.se>'
	assert_line --index 2 'Co-authored-by: B <b@x.y>'
	assert_line --index 3 'Co-authored-by: C <c@x.y>'
}

@test "git-team: (scope: global) enable 'all via -A' should provision the commit template" {
	run bash -c "/usr/local/bin/git-team enable -A &>/dev/null && cat /root/.config/git-team/commit-templates/global/COMMIT_TEMPLATE"
	assert_success
	assert_line --index 0 'Co-authored-by: A <a@x.y>'
	assert_line --index 1 'Co-authored-by: B <b@x.y>'
	assert_line --index 2 'Co-authored-by: C <c@x.y>'
}

@test "git-team: (scope: global) enable 'all via --all' should provision the commit template" {
	run bash -c "/usr/local/bin/git-team enable --all &>/dev/null && cat /root/.config/git-team/commit-templates/global/COMMIT_TEMPLATE"
	assert_success
	assert_line --index 0 'Co-authored-by: A <a@x.y>'
	assert_line --index 1 'Co-authored-by: B <b@x.y>'
	assert_line --index 2 'Co-authored-by: C <c@x.y>'
}

@test "git-team: (scope: global) enable shorthand should display the enabled co-authors in alphabetical order" {
	run /usr/local/bin/git-team b a c 'Ad-hoc <adhoc@tmp.se>'
	assert_success
	assert_line --index 0 "warn: 'git team (without further sub-command specification)' has been deprecated and is going to be removed in a future major release, use 'git team enable' instead"
	assert_line --index 1 'git-team enabled'
	assert_line --index 2 'co-authors'
	assert_line --index 3 '─ A <a@x.y>'
	assert_line --index 4 '─ Ad-hoc <adhoc@tmp.se>'
	assert_line --index 5 '─ B <b@x.y>'
	assert_line --index 6 '─ C <c@x.y>'
}

@test "git-team: (scope: global) enable should display the enabled co-authors in alphabetical order" {
	run /usr/local/bin/git-team enable b a c 'Ad-hoc <adhoc@tmp.se>'
	assert_success
	assert_line --index 0 'git-team enabled'
	assert_line --index 1 'co-authors'
	assert_line --index 2 '─ A <a@x.y>'
	assert_line --index 3 '─ Ad-hoc <adhoc@tmp.se>'
	assert_line --index 4 '─ B <b@x.y>'
	assert_line --index 5 '─ C <c@x.y>'
}

@test "git-team: (scope: global) issuing enable should be idempotent" {
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

@test "git-team: (scope: global) enable should ignore duplicates" {
	run /usr/local/bin/git-team enable a a 'A <a@x.y>'
	assert_success
	assert_line --index 0 'git-team enabled'
	assert_line --index 1 'co-authors'
	assert_line --index 2 '─ A <a@x.y>'
	assert_line --index 3 ''
}

@test "git-team: (scope: global) enable should fail when trying to enable with a non-existing alias" {
	run /usr/local/bin/git-team enable non-existing-alias
	assert_failure 255
	assert_line 'error: Failed to resolve alias team.alias.non-existing-alias'
}

