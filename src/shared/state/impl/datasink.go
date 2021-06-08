package stateimpl

import (
	"errors"

	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
	giterror "github.com/hekmekk/git-team/src/shared/gitconfig/error"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
	state "github.com/hekmekk/git-team/src/shared/state/entity"
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
func (ds GitConfigDataSink) PersistEnabled(scope activationscope.Scope, coauthors []string) error {
	return ds.persist(scope, state.NewStateEnabled(coauthors))
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

	return nil
}
