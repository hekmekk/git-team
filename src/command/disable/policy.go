package disable

import (
	"os"

	"github.com/hekmekk/git-team/src/core/events"
	giterror "github.com/hekmekk/git-team/src/core/gitconfig/error"
	commitsettingsreader "github.com/hekmekk/git-team/src/shared/commitsettings/reader"
)

// Dependencies the dependencies of the disable Policy module
type Dependencies struct {
	GitUnsetCommitTemplate func() error
	GitUnsetHooksPath      func() error
	CommitSettingsReader   commitsettingsreader.CommitSettingsReader
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

	commitTemplatePath := deps.CommitSettingsReader.Read().GitTeamCommitTemplatePath

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
