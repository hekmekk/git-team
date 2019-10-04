#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

@test "git-team: disable should persist the current status to gitconfig" {
	mkdir -p /root/.config/git-team/
	/usr/local/bin/git-team disable

	run git config --global --get-regexp team.state
	assert_success
	assert_output 'team.state.status disabled'

	rm -rf /root/.config/git-team/
}

@test "git-team: disable should disable the prepare-commit-msg hook" {
	mkdir -p /root/.config/git-team/

	run bash -c "/usr/local/bin/git-team disable &>/dev/null && git config --global core.hooksPath"
	assert_failure 1
	refute_line --regexp '\w+'

	rm -rf /root/.config/git-team/
}

@test "git-team: disable should unset the commit template" {
	mkdir -p /root/.config/git-team/

	run bash -c "/usr/local/bin/git-team disable &>/dev/null && git config --global commit.template"
	assert_failure 1
	refute_line --regexp '\w+'

	rm -rf /root/.config/git-team/
}

@test "git-team: disable should treat a previously disabled git-team idempotently" {
	mkdir -p /root/.config/git-team/

	run /usr/local/bin/git-team disable
	assert_success
	assert_line "git-team disabled"

	rm -rf /root/.config/git-team/
}

@test "git-team: disable should disable a previously enabled git-team" {
	mkdir -p /root/.config/git-team/
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>'

	run /usr/local/bin/git-team disable
	assert_success
	assert_line "git-team disabled"

	rm -rf /root/.config/git-team/
}

