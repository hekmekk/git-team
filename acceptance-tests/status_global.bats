#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

setup() {
	/usr/local/bin/git-team config activation-scope global
}

@test 'git-team: (scope: global) status should properly display a disabled status' {
	run /usr/local/bin/git-team status
	assert_success
	assert_line 'git-team disabled'
}

@test 'git-team: (scope: global) status should properly display the enabled status' {
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

