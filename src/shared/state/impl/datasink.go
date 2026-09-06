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

	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	giterror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
	state "github.com/hekmekk/git-team/v2/src/shared/state/entity"
)

// GitConfigDataSink write data directly to gitconfig
type GitConfigDataSink struct {
	GitConfigWriter gitconfig.Writer
}

// NewGitConfigDataSink construct new DataSink
func NewGitConfigDataSink(gitConfigWriter gitconfig.Writer) GitConfigDataSink {
	return GitConfigDataSink{GitConfigWriter: gitConfigWriter}
}

// PersistEnabled persist the current state as enabled
func (ds GitConfigDataSink) PersistEnabled(scope activationscope.Scope, coauthors []string, previousHooksPath string) error {
	return ds.persist(scope, state.NewStateEnabled(coauthors, previousHooksPath))
}

// PersistDisabled persist the current state as disabled
func (ds GitConfigDataSink) PersistDisabled(scope activationscope.Scope) error {
	return ds.persist(scope, state.NewStateDisabled())
}

func (ds GitConfigDataSink) persist(activationScope activationscope.Scope, state state.State) error {
	gitConfigWriter := ds.GitConfigWriter

	var gitConfigScope gitconfigscope.Scope
	if activationScope == activationscope.Global {
		gitConfigScope = gitconfigscope.Global
	} else {
		gitConfigScope = gitconfigscope.Local
	}

	if err := gitConfigWriter.UnsetAll(gitConfigScope, "team.state.active-coauthors"); err != nil && !errors.Is(err, giterror.ErrTryingToUnsetAnOptionWhichDoesNotExist) {
		return errors.New("failed to unset team.state.active-coauthors")
	}

	for _, coauthor := range state.Coauthors {
		if err := gitConfigWriter.Add(gitConfigScope, "team.state.active-coauthors", coauthor); err != nil {
			return errors.New("failed to set team.state.active-coauthors")
		}
	}

	if err := gitConfigWriter.ReplaceAll(gitConfigScope, "team.state.status", string(state.Status)); err != nil {
		return errors.New("failed to replace team.state.status")
	}

	if state.PreviousHooksPath == "" {
		if err := gitConfigWriter.UnsetAll(gitConfigScope, "team.state.previous-hooks-path"); err != nil && !errors.Is(err, giterror.ErrTryingToUnsetAnOptionWhichDoesNotExist) {
			return errors.New("failed to unset team.state.previous-hooks-path")
		}
	} else {
		if err := gitConfigWriter.ReplaceAll(gitConfigScope, "team.state.previous-hooks-path", state.PreviousHooksPath); err != nil {
			return errors.New("failed to replace team.state.previous-hooks-path")
		}
	}

	return nil
}
