package removecmdadapter

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/v2/src/command/assignments/remove"
	removeeventadapter "github.com/hekmekk/git-team/v2/src/command/assignments/remove/cliadapter/event"
	commandadapter "github.com/hekmekk/git-team/v2/src/shared/cli/commandadapter"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
	aliascompletion "github.com/hekmekk/git-team/v2/src/shared/completion"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl"
	gitconfiglegacy "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl/legacy"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
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
