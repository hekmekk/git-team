package listcmdadapter

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
	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/v2/src/command/assignments/list"
	listeventadapter "github.com/hekmekk/git-team/v2/src/command/assignments/list/cliadapter/event"
	commandadapter "github.com/hekmekk/git-team/v2/src/shared/cli/commandadapter"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl"
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
