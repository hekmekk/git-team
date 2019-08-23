package removecmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/core/gitconfig"
	"github.com/hekmekk/git-team/src/remove"
)

// Definition the command, arguments, and dependencies
type Definition struct {
	CommandName string
	Policy      remove.Policy
}

// NewDefinition the constructor for Definition
func NewDefinition(app *kingpin.Application) Definition {

	command := app.Command("rm", "Remove an alias to co-author assignment")
	alias := command.Arg("alias", "The alias to identify the assignment to be removed").Required().String()

	return Definition{
		CommandName: command.FullCommand(),
		Policy: remove.Policy{
			Req: remove.DeAllocationRequest{
				Alias: alias,
			},
			Deps: remove.Dependencies{
				GitRemoveAlias: gitconfig.RemoveAlias,
			},
		},
	}
}
