#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

setup() {
	git config --global init.defaultBranch main

	/usr/local/bin/git-team config activation-scope global
}

teardown() {
	rm /home/git-team-acceptance-test/.gitconfig
}

@test "git-team: (scope: global) disable should disable a previously enabled git-team" {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'

	run /usr/local/bin/git-team disable
	assert_success
	assert_line "git-team disabled"
}

@test "git-team: (scope: global) disable should persist the current status to gitconfig" {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run git config --global --get-regexp team.state
	assert_success
	assert_output 'team.state.status disabled'
}

@test "git-team: (scope: global) disable should persist a previous hooks path as the current hooks path" {
	git config --global core.hooksPath "/path/to/non-git-team-hooks"
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "git config --global core.hooksPath"
	assert_success
	assert_line --index 0 '/path/to/non-git-team-hooks'

	git config --global --unset core.hooksPath | true
}

@test "git-team: (scope: global) disable should disable the prepare-commit-msg hook" {
	/usr/local/bin/git-team config activation-scope global
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "git config --global core.hooksPath"
	assert_failure 1
	refute_line --regexp '\w+'
}

@test "git-team: (scope: global) disable should unset the commit template" {
	/usr/local/bin/git-team config activation-scope global
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "git config --global commit.template"
	assert_failure 1
	refute_line --regexp '\w+'
}

@test "git-team: (scope: global) disable should remove the according COMMIT_TEMPLATE" {
	/usr/local/bin/git-team config activation-scope global
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'
	/usr/local/bin/git-team disable

	run bash -c "ls -la /home/git-team-acceptance-test/.git-team/commit-templates/global/COMMIT_TEMPLATE"
	assert_failure 1
	assert_line "ls: /home/git-team-acceptance-test/.git-team/commit-templates/global/COMMIT_TEMPLATE: No such file or directory"
}

@test "git-team: (scope: global) disable should treat a previously disabled git-team in an idempotent way" {
	run /usr/local/bin/git-team disable
	assert_success
	assert_line "git-team disabled"
}

