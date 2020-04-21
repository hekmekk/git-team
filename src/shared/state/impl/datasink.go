package stateimpl

import (
	giterror "github.com/hekmekk/git-team/src/shared/gitconfig/error"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
	scope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
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
func (ds GitConfigDataSink) PersistEnabled(scope scope.Scope, coauthors []string) error {
	return ds.persist(scope, state.NewStateEnabled(coauthors))
}

// PersistDisabled persist the current state as disabled
func (ds GitConfigDataSink) PersistDisabled(scope scope.Scope) error {
	return ds.persist(scope, state.NewStateDisabled())
}

func (ds GitConfigDataSink) persist(scope scope.Scope, state state.State) error {
	gitConfigWriter := ds.GitConfigWriter

	if err := gitConfigWriter.UnsetAll(scope, "team.state.active-coauthors"); err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return err
	}

	for _, coauthor := range state.Coauthors {
		if err := gitConfigWriter.Add(scope, "team.state.active-coauthors", coauthor); err != nil {
			return err
		}
	}

	if err := gitConfigWriter.ReplaceAll(scope, "team.state.status", string(state.Status)); err != nil {
		return err
	}

	return nil
}
