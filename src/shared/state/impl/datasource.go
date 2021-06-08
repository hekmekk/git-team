package stateimpl

import (
	"fmt"

	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
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

	return state.NewStateEnabled(activeCoauthors), nil
}
