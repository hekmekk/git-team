package disable

import (
	"os"

	"github.com/hekmekk/git-team/src/core/events"
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/shared/config/interface"
	giterror "github.com/hekmekk/git-team/src/shared/gitconfig/error"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
	state "github.com/hekmekk/git-team/src/shared/state/interface"
)

// Dependencies the dependencies of the disable Policy module
type Dependencies struct {
	GitConfigReader gitconfig.Reader
	GitConfigWriter gitconfig.Writer
	StatFile        func(string) (os.FileInfo, error)
	RemoveFile      func(string) error
	StateWriter     state.Writer
	ConfigReader    config.Reader
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply disable team mode
func (policy Policy) Apply() events.Event {
	deps := policy.Deps
	gitConfigWriter := deps.GitConfigWriter

	if err := gitConfigWriter.UnsetAll(gitconfigscope.Global, "core.hooksPath"); err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return Failed{Reason: err}
	}

	commitTemplatePath, err := deps.GitConfigReader.Get(gitconfigscope.Global, "commit.template")
	if err != nil && err.Error() != giterror.SectionOrKeyIsInvalid {
		return Failed{Reason: err}
	}

	if err := gitConfigWriter.UnsetAll(gitconfigscope.Global, "commit.template"); err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return Failed{Reason: err}
	}

	if _, err := deps.StatFile(commitTemplatePath); err == nil {
		if err := deps.RemoveFile(commitTemplatePath); err != nil {
			return Failed{Reason: err}
		}
	}

	if err := deps.StateWriter.PersistDisabled(activationscope.Global); err != nil {
		return Failed{Reason: err}
	}

	return Succeeded{}
}
