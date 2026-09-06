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

@test 'git-team: completion should show the available scripts' {
	run /usr/local/bin/git-team completion

	assert_success
	assert_line --index 6 'COMMANDS:'
	assert_line --index 7 '   bash     Bash completion'
	assert_line --index 8 '   zsh      Zsh completion'
}

@test 'git-team: completion bash should print the bash completion script' {
	run /usr/local/bin/git-team completion bash

	assert_success
	assert_line --index 0 '#!/bin/bash'
	assert_line --index 14 '_git_team() {'
	assert_line --index 27 '}'
	assert_line --index 29 '_git_team_bash_completion() {'
	assert_line --index 42 '}'
	assert_line --index 43 'complete -F _git_team_bash_completion git-team'
}

@test 'git-team: completion zsh should print the zsh completion script' {
	run /usr/local/bin/git-team completion zsh

	assert_success
	assert_line --index 0 '#compdef git-team'
	assert_line --index 13 'function _git-team {'
	assert_line --index 35 '}'
	assert_line --index 36 'compdef _git-team git-team'
	assert_line --index 39 "zstyle ':completion:*:*:git:*' user-commands team:'manage and enhance git commit messages with co-authors'"
}

