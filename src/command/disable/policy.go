package disable

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hekmekk/git-team/src/core/events"
	activation "github.com/hekmekk/git-team/src/shared/activation/interface"
	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/src/shared/config/interface"
	giterror "github.com/hekmekk/git-team/src/shared/gitconfig/error"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
	state "github.com/hekmekk/git-team/src/shared/state/interface"
)

// Dependencies the dependencies of the disable Policy module
type Dependencies struct {
	GitConfigReader     gitconfig.Reader
	GitConfigWriter     gitconfig.Writer
	StatFile            func(string) (os.FileInfo, error)
	RemoveFile          func(string) error
	StateReader         state.Reader
	StateWriter         state.Writer
	ConfigReader        config.Reader
	ActivationValidator activation.Validator
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply disable team mode
func (policy Policy) Apply() events.Event {
	deps := policy.Deps
	gitConfigWriter := deps.GitConfigWriter

	cfg, err := deps.ConfigReader.Read()
	if err != nil {
		return Failed{Reason: fmt.Errorf("failed to read config: %s", err)}
	}

	activationScope := cfg.ActivationScope

	if activationScope == activationscope.RepoLocal && !deps.ActivationValidator.IsInsideAGitRepository() {
		return Failed{Reason: fmt.Errorf("failed to disable with activation-scope=%s: not inside a git repository", activationScope)}
	}

	var gitConfigScope gitconfigscope.Scope
	if activationScope == activationscope.Global {
		gitConfigScope = gitconfigscope.Global
	} else {
		gitConfigScope = gitconfigscope.Local
	}

	appState, err := deps.StateReader.Query(activationScope)
	if err != nil {
		return Failed{Reason: fmt.Errorf("failed to read current state: %s", err)}
	}

	previousHooksPath := appState.PreviousHooksPath
	if "" == previousHooksPath {
		if err := gitConfigWriter.UnsetAll(gitConfigScope, "core.hooksPath"); err != nil && !errors.Is(err, giterror.ErrTryingToUnsetAnOptionWhichDoesNotExist) {
			return Failed{Reason: fmt.Errorf("failed to unset core.hooksPath: %s", err)}
		}
	} else {
		if err := gitConfigWriter.ReplaceAll(gitConfigScope, "core.hooksPath", previousHooksPath); err != nil {
			return Failed{Reason: fmt.Errorf("failed to replace core.hooksPath: %s", err)}
		}
	}

	commitTemplatePath, err := deps.GitConfigReader.Get(gitConfigScope, "commit.template")
	if err != nil && !errors.Is(err, giterror.ErrSectionOrKeyIsInvalid) {
		return Failed{Reason: fmt.Errorf("failed to get commit.template: %s", err)}
	}

	if err := gitConfigWriter.UnsetAll(gitConfigScope, "commit.template"); err != nil && !errors.Is(err, giterror.ErrTryingToUnsetAnOptionWhichDoesNotExist) {
		return Failed{Reason: fmt.Errorf("failed to unset commit.template: %s", err)}
	}

	if commitTemplatePath != "" {
		var templatePathToDelete string
		if activationScope == activationscope.Global {
			templatePathToDelete = commitTemplatePath
		} else {
			templatePathToDelete = filepath.Dir(commitTemplatePath)
		}

		if _, err := deps.StatFile(templatePathToDelete); err == nil {
			if err := deps.RemoveFile(templatePathToDelete); err != nil {
				return Failed{Reason: fmt.Errorf("failed to remove commit template: %s", err)}
			}
		}
	}

	if err := deps.StateWriter.PersistDisabled(activationScope); err != nil {
		return Failed{Reason: fmt.Errorf("failed to write current state: %s", err)}
	}

	return Succeeded{}
}
