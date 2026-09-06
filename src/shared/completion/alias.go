package completion

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"sort"
	"strings"

	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// AliasShellCompletion generate completion
type AliasShellCompletion struct {
	GitConfigReader gitconfig.Reader
}

// NewAliasShellCompletion construct new CoAuthorShellCompletion
func NewAliasShellCompletion(gitconfigReader gitconfig.Reader) AliasShellCompletion {
	return AliasShellCompletion{
		GitConfigReader: gitconfigReader,
	}
}

// Complete return not yet selected aliases
func (completion AliasShellCompletion) Complete(selectedAliases []string) []string {
	allAssignments, err := completion.GitConfigReader.GetRegexp(gitconfigscope.Global, "team.alias")

	if err != nil {
		return []string{}
	}

	remainingAliases := []string{}

	for rawAlias := range allAssignments {
		alias := strings.TrimPrefix(rawAlias, "team.alias.")
		isSelected := false
		for _, selectedAlias := range selectedAliases {
			if selectedAlias == alias {
				isSelected = true
			}
		}

		if !isSelected {
			remainingAliases = append(remainingAliases, alias)
		}
	}

	sort.Strings(remainingAliases)

	return remainingAliases
}
