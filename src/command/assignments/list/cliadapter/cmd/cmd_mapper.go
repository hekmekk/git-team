package listcmdadapter

import (
	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/src/command/assignments/list"
	listeventadapter "github.com/hekmekk/git-team/src/command/assignments/list/cliadapter/event"
	commandadapter "github.com/hekmekk/git-team/src/shared/cli/commandadapter"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
)

// Command the ls command
func Command() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List your assignments",
		Action: func(c *cli.Context) error {
			return commandadapter.Run(policy(), listeventadapter.MapEventToEffect)
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
