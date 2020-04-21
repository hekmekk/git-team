package stateimpl

import (
	"errors"

	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
	state "github.com/hekmekk/git-team/src/shared/state/entity"
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
func (ds GitConfigDataSource) Query(scope activationscope.ActivationScope) (state.State, error) {
	status, err := ds.GitConfigReader.Get(gitconfigscope.Global, "team.state.status")
	if err != nil || "disabled" == status || "" == status {
		return state.NewStateDisabled(), nil
	}

	activeCoauthors, err := ds.GitConfigReader.GetAll(gitconfigscope.Global, "team.state.active-coauthors")
	if err != nil {
		return state.State{}, errors.New("no active co-authors found")
	}

	return state.NewStateEnabled(activeCoauthors), nil
}
