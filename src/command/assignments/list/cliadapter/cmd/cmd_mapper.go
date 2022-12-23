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
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "only-alias", Value: false, Aliases: []string{"o"}, Usage: "Only show the alias, omitting the leading dash and the associated co-authors"},
		},
		Action: func(c *cli.Context) error {
			onlyAlias := c.Bool("only-alias")
			return commandadapter.Run(policy(&onlyAlias), listeventadapter.MapEventToEffect)
		},
	}
}

func policy(onlyAlias *bool) list.Policy {
	return list.Policy{
		Req: list.ListRequest{
			OnlyAlias: onlyAlias,
		},
		Deps: list.Dependencies{
			GitConfigReader: gitconfig.NewDataSource(),
		},
	}
}
