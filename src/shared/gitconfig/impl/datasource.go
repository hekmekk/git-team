package gitconfigimpl

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
	gitconfiglegacy "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl/legacy"
	scope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// DataSource read data directly from gitconfig
type DataSource struct{}

// NewDataSource construct new DataSource
func NewDataSource() DataSource {
	return DataSource{}
}

// Get read the first value for a key
func (ds DataSource) Get(scope scope.Scope, key string) (string, error) {
	return gitconfiglegacy.Get(scope, key)
}

// GetAll read all values for a key
func (ds DataSource) GetAll(scope scope.Scope, key string) ([]string, error) {
	return gitconfiglegacy.GetAll(scope, key)
}

// GetRegexp read all values matching a pattern
func (ds DataSource) GetRegexp(scope scope.Scope, pattern string) (map[string]string, error) {
	return gitconfiglegacy.GetRegexp(scope, pattern)
}

// List show the entire config
func (ds DataSource) List(scope scope.Scope) (map[string]string, error) {
	return gitconfiglegacy.List(scope)
}
