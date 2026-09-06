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

// DataSink write data directly to gitconfig
type DataSink struct{}

// NewDataSink construct new DataSink
func NewDataSink() DataSink {
	return DataSink{}
}

// Add add a value to a key
func (ds DataSink) Add(scope scope.Scope, key string, value string) error {
	return gitconfiglegacy.Add(scope, key, value)
}

// ReplaceAll modify a setting
func (ds DataSink) ReplaceAll(scope scope.Scope, key string, value string) error {
	return gitconfiglegacy.ReplaceAll(scope, key, value)
}

// UnsetAll remove a setting
func (ds DataSink) UnsetAll(scope scope.Scope, key string) error {
	return gitconfiglegacy.UnsetAll(scope, key)
}
