package disablecmdadapter

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/command/disable"
	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/gitconfig"
	"github.com/hekmekk/git-team/src/core/state_repository"
)

// Definition definition of the disable command
type Definition struct {
	CommandName string
	Policy      disable.Policy
}

// NewDefinition the constructor for Definition
func NewDefinition(app *kingpin.Application) Definition {

	command := app.Command("disable", "Use default commit template and remove prepare-commit-msg hook")

	return Definition{
		CommandName: command.FullCommand(),
		Policy: disable.Policy{
			Deps: disable.Dependencies{
				GitUnsetCommitTemplate: gitconfig.UnsetCommitTemplate,
				GitUnsetHooksPath:      func() error { return gitconfig.UnsetAll("core.hooksPath") },
				LoadConfig:             config.Load,
				StatFile:               os.Stat,
				RemoveFile:             os.Remove,
				PersistDisabled:        staterepository.PersistDisabled,
			},
		},
	}
}
