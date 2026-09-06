package datasink

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
	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// GitconfigDataSink writes configuration to gitconfig
type GitconfigDataSink struct {
	GitConfigWriter gitconfig.Writer
}

// NewGitconfigDataSink constructs new GitconfigDataSink
func NewGitconfigDataSink(gitConfigWriter gitconfig.Writer) GitconfigDataSink {
	return GitconfigDataSink{GitConfigWriter: gitConfigWriter}
}

// SetActivationScope write activation-scope setting to gitconfig
func (ds GitconfigDataSink) SetActivationScope(scope activationscope.Scope) error {
	return ds.GitConfigWriter.ReplaceAll(gitconfigscope.Global, "team.config.activation-scope", scope.String())
}
