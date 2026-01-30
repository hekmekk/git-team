package stateimpl

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
