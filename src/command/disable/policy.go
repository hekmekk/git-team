package disable

import (
	"os"

	"github.com/hekmekk/git-team/src/core/events"
	gitconfig "github.com/hekmekk/git-team/src/core/gitconfig"
	giterror "github.com/hekmekk/git-team/src/core/gitconfig/error"
)

// Dependencies the dependencies of the disable Policy module
type Dependencies struct {
	GitSettingsReader gitconfig.SettingsReader
	GitSettingsWriter gitconfig.SettingsWriter
	StatFile          func(string) (os.FileInfo, error)
	RemoveFile        func(string) error
	PersistDisabled   func() error
	// TODO: ConfigReader config.Reader to get activation scope from config
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply disable team mode
func (policy Policy) Apply() events.Event {
	deps := policy.Deps
	gitSettingsWriter := deps.GitSettingsWriter

	if err := gitSettingsWriter.Unset("core.hooksPath"); err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return Failed{Reason: err}
	}

	commitTemplatePath, err := deps.GitSettingsReader.Get("commit.template")
	if err != nil && err.Error() != giterror.SectionOrKeyIsInvalid {
		return Failed{Reason: err}
	}

	if err := gitSettingsWriter.Unset("commit.template"); err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return Failed{Reason: err}
	}

	if _, err := deps.StatFile(commitTemplatePath); err == nil {
		if err := deps.RemoveFile(commitTemplatePath); err != nil {
			return Failed{Reason: err}
		}
	}

	if err := deps.PersistDisabled(); err != nil {
		return Failed{Reason: err}
	}

	return Succeeded{}
}
