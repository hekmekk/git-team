package statuscmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/core/state_repository"
	"github.com/hekmekk/git-team/src/status"
)

// Definition the command, arguments, and dependencies
type Definition struct {
	CommandName string
	Policy      status.Policy
}

// NewDefinition the constructor for Definition
func NewDefinition(app *kingpin.Application) Definition {

	command := app.Command("status", "Print the current status")

	return Definition{
		CommandName: command.FullCommand(),
		Policy: status.Policy{
			Deps: status.Dependencies{
				StateRepositoryQuery: staterepository.Query,
			},
		},
	}
}
