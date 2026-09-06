package commandadapter

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

// TODO: this should live somewhere else...
// TODO: this should depend on the gitconfig interface only as well

import (
	"fmt"

	gitconfiglegacy "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl/legacy"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// ResolveAliases convenience function to resolve multiple aliases and accumulate errors
func ResolveAliases(aliases []string) ([]string, []error) {
	return resolveAliases(ResolveAlias)(aliases)
}

func resolveAliases(resolveAlias func(string) (string, error)) func([]string) ([]string, []error) {
	return func(aliases []string) ([]string, []error) {
		var resolvedAliases []string
		var resolveErrors []error

		for _, alias := range aliases {
			var resolvedCoauthor, err = resolveAlias(alias)
			if err != nil {
				resolveErrors = append(resolveErrors, err)
			} else {
				resolvedAliases = append(resolvedAliases, resolvedCoauthor)
			}
		}

		return resolvedAliases, resolveErrors
	}
}

// ResolveAlias lookup "team.alias.<alias>" globally
func ResolveAlias(alias string) (string, error) {
	return resolveAlias(gitconfiglegacy.Get)(alias)
}

func resolveAlias(gitconfigGet func(gitconfigscope.Scope, string) (string, error)) func(string) (string, error) {
	return func(alias string) (string, error) {
		coauthor, err := gitconfigGet(gitconfigscope.Global, fmt.Sprintf("team.alias.%s", alias))
		if err != nil || coauthor == "" {
			return "", fmt.Errorf("failed to resolve alias team.alias.%s", alias)
		}

		return coauthor, nil
	}
}
