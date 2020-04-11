#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

@test "git-team: config should show the current configuration" {
	run bash -c "git team config"
	assert_success
	assert_line --index 0 'config'
	assert_line --index 1 '─ [ro] commit-template-path: /root/.config/git-team/COMMIT_TEMPLATE'
	assert_line --index 2 '─ [ro] hooks-path: /usr/local/etc/git-team/hooks'
	assert_line --index 3 '─ [rw] activation-scope: global'
}

