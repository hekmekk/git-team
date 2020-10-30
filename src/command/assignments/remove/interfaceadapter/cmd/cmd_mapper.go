package removecmdadapter

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/assignments/remove"
	removeeventadapter "github.com/hekmekk/git-team/src/command/assignments/remove/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/effects"
	aliascompletion "github.com/hekmekk/git-team/src/shared/completion"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
	gitconfiglegacy "github.com/hekmekk/git-team/src/shared/gitconfig/impl/legacy"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

// Command the rm command
func Command() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Aliases:   []string{"rm"},
		Usage:     "Remove an alias to co-author assignment",
		ArgsUsage: "<alias>",
		Action: func(c *cli.Context) error {
			args := c.Args()
			if args.Len() != 1 {
				effects.NewPrintErr(errors.New("exactly one alias must be specified")).Run()
				return nil
			}

			alias := args.First()
			return commandadapter.RunUrFave(policy(&alias), removeeventadapter.MapEventToEffects)(c)
		},
		BashComplete: func(c *cli.Context) {
			args := c.Args()
			if args.Len() == 0 {
				remainingAliases := aliascompletion.NewAliasShellCompletion(gitconfig.NewDataSource()).Complete(args.Slice())
				for _, alias := range remainingAliases {
					fmt.Println(alias)
				}
			} else {
				fmt.Println()
			}
		},
	}
}

func policy(alias *string) remove.Policy {
	return remove.Policy{
		Req: remove.DeAllocationRequest{
			Alias: alias,
		},
		Deps: remove.Dependencies{
			GitRemoveAlias: func(alias string) error {
				return gitconfiglegacy.UnsetAll(gitconfigscope.Global, fmt.Sprintf("team.alias.%s", alias))
			},
		},
	}
}
