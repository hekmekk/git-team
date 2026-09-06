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

	/usr/local/bin/git-team disable
	touch /tmp/COMMIT_MSG
}

teardown() {
	rm /tmp/COMMIT_MSG
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - message" {
	run bash -c "/usr/local/bin/prepare-commit-msg-git-team /tmp/COMMIT_MSG message && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - none" {
	run bash -c "/usr/local/bin/prepare-commit-msg-git-team /tmp/COMMIT_MSG && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - commit" {
	run bash -c "/usr/local/bin/prepare-commit-msg-git-team /tmp/COMMIT_MSG commit && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - template" {
	run bash -c "/usr/local/bin/prepare-commit-msg-git-team /tmp/COMMIT_MSG template && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - merge" {
	run bash -c "/usr/local/bin/prepare-commit-msg-git-team /tmp/COMMIT_MSG merge && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

@test "prepare-commit-msg: git-team disabled: (scope: global) - squash" {
	run bash -c "/usr/local/bin/prepare-commit-msg-git-team /tmp/COMMIT_MSG squash && cat /tmp/COMMIT_MSG"
	assert_success
	refute_output --regexp '\w+'
}

