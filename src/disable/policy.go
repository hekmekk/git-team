package disable

import (
	"fmt"
	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/events"
	giterror "github.com/hekmekk/git-team/src/core/gitconfig/error"
	"os"
)

// Dependencies the dependencies of the disable Policy module
type Dependencies struct {
	GitUnsetCommitTemplate func() error
	GitUnsetHooksPath      func() error
	LoadConfig             func() (config.Config, error)
	StatFile               func(string) (os.FileInfo, error)
	RemoveFile             func(string) error
	PersistDisabled        func() error
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply disable team mode
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	if err := deps.GitUnsetHooksPath(); err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return Failed{Reason: err}
	}

	if err := deps.GitUnsetCommitTemplate(); err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return Failed{Reason: err}
	}

	cfg, err := deps.LoadConfig()
	if err != nil {
		return Failed{Reason: err}
	}

	templateFilePath := fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName)

	if _, err := deps.StatFile(templateFilePath); err == nil {
		if err := deps.RemoveFile(templateFilePath); err != nil {
			return Failed{Reason: err}
		}
	}

	if err := deps.PersistDisabled(); err != nil {
		return Failed{Reason: err}
	}

	return Succeeded{}
}