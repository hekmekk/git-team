#!/usr/bin/env bats

# Copyright (C) 2026 Rea Sand

# This file is part of git-team.

# This program is free software: you can redistribute it and/or
# modify it under the terms of the GNU General Public License
# as published by the Free Software Foundation, either version 3
# of the License, or (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program. If not, see <https://www.gnu.org/licenses/>.

setup() {
	bats_load_library bats-support
	bats_load_library bats-assert

	/usr/local/bin/git-team config activation-scope global
}

@test 'git-team: (scope: global) status should properly display the disabled status' {
	run /usr/local/bin/git-team status
	assert_success
	assert_line 'git-team disabled'
}

@test 'git-team: (scope: global) status should properly display the disabled status in a json format' {
	run /usr/local/bin/git-team status --json
	assert_success
	assert_line --index 0 '{"status":"disabled","coAuthors":[]}'
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

@test 'git-team: (scope: global) status should properly display the enabled status in a json format' {
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'

	run /usr/local/bin/git-team status --json
	assert_success
	assert_line --index 0 '{"status":"enabled","coAuthors":["A <a@x.y>","B <b@x.y>","C <c@x.y>"]}'

	/usr/local/bin/git-team disable
}

