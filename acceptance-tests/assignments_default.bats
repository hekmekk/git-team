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
}

@test "git-team: assignments (default) should show 'no assignments'" {
	run /usr/local/bin/git-team assignments
	assert_success
	assert_line --index 0 'No assignments'
}

@test "git-team: assignments (default) should show all alias -> coauthor assignments" {
	/usr/local/bin/git-team assignments add a 'A <a@x.y>'
	/usr/local/bin/git-team assignments add bb 'B <b@x.y>'
	/usr/local/bin/git-team assignments add c 'C <c@x.y>'

	run /usr/local/bin/git-team assignments
	assert_success
	assert_line --index 0 'Assignments'
	assert_line --index 1 '─ a  →  A <a@x.y>'
	assert_line --index 2 '─ bb →  B <b@x.y>'
	assert_line --index 3 '─ c  →  C <c@x.y>'

	/usr/local/bin/git-team assignments rm a
	/usr/local/bin/git-team assignments rm bb
	/usr/local/bin/git-team assignments rm c
}

@test "git-team: assignments ls should show all alias -> coauthor assignments" {
	/usr/local/bin/git-team assignments add a 'A <a@x.y>'
	/usr/local/bin/git-team assignments add bb 'B <b@x.y>'
	/usr/local/bin/git-team assignments add c 'C <c@x.y>'

	run /usr/local/bin/git-team assignments ls
	assert_success
	assert_line --index 0 'Assignments'
	assert_line --index 1 '─ a  →  A <a@x.y>'
	assert_line --index 2 '─ bb →  B <b@x.y>'
	assert_line --index 3 '─ c  →  C <c@x.y>'

	/usr/local/bin/git-team assignments rm a
	/usr/local/bin/git-team assignments rm bb
	/usr/local/bin/git-team assignments rm c
}

