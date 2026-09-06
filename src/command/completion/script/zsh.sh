#compdef git-team

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

function _git-team {
  local -a opts
  local cmd cur
  cmd=${words[1]}
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    if [[ "${cmd}" == "team" ]]; then
      opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 git ${cmd} ${words[@]:1:#words[@]-1} ${cur} --generate-bash-completion)}")
    else
      opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${cmd} ${words[@]:1:#words[@]-1} ${cur} --generate-bash-completion)}")
    fi
  else
    if [[ "${cmd}" == "team" ]]; then
      opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 git ${cmd} ${words[@]:1:#words[@]-1} --generate-bash-completion)}")
    else
      opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${cmd} ${words[@]:1:#words[@]-1} --generate-bash-completion)}")
    fi
  fi
  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  fi
  return
}
compdef _git-team git-team
# the function will automatically be discovered by the git zsh completion as its name follows the pattern _git-<user-command>
# make the 'team' user-command known to the git completion, so that it will be suggested using git<TAB>
zstyle ':completion:*:*:git:*' user-commands team:'manage and enhance git commit messages with co-authors'
