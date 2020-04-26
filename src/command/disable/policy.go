package disable

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hekmekk/git-team/src/core/events"
	activation "github.com/hekmekk/git-team/src/shared/activation/interface"
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
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
		return Failed{Reason: err}
	}

	activationScope := cfg.ActivationScope

	if activationScope == activationscope.RepoLocal && !deps.ActivationValidator.IsInsideAGitRepository() {
		return Failed{Reason: fmt.Errorf("Failed to disable with activation-scope=%s: not inside a git repository", activationScope)}
	}

	var gitConfigScope gitconfigscope.Scope
	if activationScope == activationscope.Global {
		gitConfigScope = gitconfigscope.Global
	} else {
		gitConfigScope = gitconfigscope.Local
	}

	if err := gitConfigWriter.UnsetAll(gitConfigScope, "core.hooksPath"); err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return Failed{Reason: err}
	}

	commitTemplatePath, err := deps.GitConfigReader.Get(gitConfigScope, "commit.template")
	if err != nil && err.Error() != giterror.SectionOrKeyIsInvalid {
		return Failed{Reason: err}
	}

	if err := gitConfigWriter.UnsetAll(gitConfigScope, "commit.template"); err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return Failed{Reason: err}
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
				return Failed{Reason: err}
			}
		}
	}

	if err := deps.StateWriter.PersistDisabled(activationScope); err != nil {
		return Failed{Reason: err}
	}

	return Succeeded{}
}
