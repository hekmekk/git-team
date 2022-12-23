package removecmdadapter

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/src/command/assignments/remove"
	removeeventadapter "github.com/hekmekk/git-team/src/command/assignments/remove/cliadapter/event"
	commandadapter "github.com/hekmekk/git-team/src/shared/cli/commandadapter"
	"github.com/hekmekk/git-team/src/shared/cli/effects"
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
				return effects.NewExitErrMsg(errors.New("exactly one alias must be specified")).Run()
			}

			alias := args.First()
			return commandadapter.Run(policy(&alias), removeeventadapter.MapEventToEffect)
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
