package stateimpl

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
	"errors"
	"fmt"
	giterror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"

	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
	state "github.com/hekmekk/git-team/v2/src/shared/state/entity"
)

// GitConfigDataSource the data source for the state reader
type GitConfigDataSource struct {
	GitConfigReader gitconfig.Reader
}

// NewGitConfigDataSource construct a new GitConfigDataSource
func NewGitConfigDataSource(gitConfigReader gitconfig.Reader) GitConfigDataSource {
	return GitConfigDataSource{GitConfigReader: gitConfigReader}
}

// Query read the current state from gitconfig
func (ds GitConfigDataSource) Query(activationScope activationscope.Scope) (state.State, error) {
	var gitConfigScope gitconfigscope.Scope
	if activationScope == activationscope.Global {
		gitConfigScope = gitconfigscope.Global
	} else {
		gitConfigScope = gitconfigscope.Local
	}

	status, err := ds.GitConfigReader.Get(gitConfigScope, "team.state.status")
	if err != nil || "disabled" == status || "" == status {
		return state.NewStateDisabled(), nil
	}

	activeCoauthors, err := ds.GitConfigReader.GetAll(gitConfigScope, "team.state.active-coauthors")
	if err != nil {
		return state.State{}, fmt.Errorf("no active co-authors found: %s", err)
	}

	previousHooksPath, err := ds.GitConfigReader.Get(gitConfigScope, "team.state.previous-hooks-path")
	if err != nil && !errors.Is(err, giterror.ErrSectionOrKeyIsInvalid) {
		return state.State{}, fmt.Errorf("failed to get previous hooks path: %s", err)
	}

	return state.NewStateEnabled(activeCoauthors, previousHooksPath), nil
}
