package stateimpl

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
