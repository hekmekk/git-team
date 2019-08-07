#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

@test 'git-team: status should properly display a disabled status' {
	run /usr/local/bin/git-team status
	assert_success
	assert_line 'git-team disabled.'
}

@test 'git-team: status should properly disaplay the enabled status' {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'

	run /usr/local/bin/git-team status
	assert_success
	assert_line --index 0 'git-team enabled.'
	assert_line --index 1 'Co-authors:'
	assert_line --index 2 '-----------'
	assert_line --index 3 'A <a@x.y>'
	assert_line --index 4 'B <b@x.y>'
	assert_line --index 5 'C <c@x.y>'

	/usr/local/bin/git-team disable
}

