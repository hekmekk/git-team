package listcmdadapter

import (
	"github.com/urfave/cli/v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/assignments/list"
	listeventadapter "github.com/hekmekk/git-team/src/command/assignments/list/interfaceadapter/event"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
)

// Command the ls command
func Command() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List your assignments",
		Action: func(c *cli.Context) error {
			return commandadapter.RunUrFave(policy(), listeventadapter.MapEventToEffects)(c)
		},
	}
}

func policy() list.Policy {
	return list.Policy{
		Deps: list.Dependencies{
			GitConfigReader: gitconfig.NewDataSource(),
		},
	}
}
