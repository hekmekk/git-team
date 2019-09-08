package listcmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/core/gitconfig"
	"github.com/hekmekk/git-team/src/list"
)

// Definition definition of the list command
type Definition struct {
	CommandName string
	Policy      list.Policy
}

// NewDefinition the constructor for Definition
func NewDefinition(app *kingpin.Application) Definition {
	command := app.Command("ls", "List currently available aliases")
	command.Alias("list")

	return Definition{
		CommandName: command.FullCommand(),
		Policy: list.Policy{
			Deps: list.Dependencies{
				GitGetAssignments: func() (map[string]string, error) { return gitconfig.GetRegexp("team.alias") },
			},
		},
	}
}
