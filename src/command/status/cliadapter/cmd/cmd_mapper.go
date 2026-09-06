package statuscmdadapter

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

	"github.com/hekmekk/git-team/v2/src/command/status"
	statuseventadapter "github.com/hekmekk/git-team/v2/src/command/status/cliadapter/event"
	activation "github.com/hekmekk/git-team/v2/src/shared/activation/impl"
	commandadapter "github.com/hekmekk/git-team/v2/src/shared/cli/commandadapter"
	config "github.com/hekmekk/git-team/v2/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl"
	state "github.com/hekmekk/git-team/v2/src/shared/state/impl"
)

// Command the status command
func Command() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Print the current status",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "json", Value: false, Usage: "Display the current status as a JSON formatted struct"},
		},
		Action: func(c *cli.Context) error {
			stateAsJson := c.Bool("json")

			return commandadapter.Run(Policy(stateAsJson), statuseventadapter.MapEventToEffect)
		},
	}
}

// Policy the status policy constructor
func Policy(stateAsJson bool) status.Policy {
	return status.Policy{
		Deps: status.Dependencies{
			ConfigReader:        config.NewGitconfigDataSource(gitconfig.NewDataSource()),
			StateReader:         state.NewGitConfigDataSource(gitconfig.NewDataSource()),
			ActivationValidator: activation.NewGitConfigDataSource(gitconfig.NewDataSource()),
			StateAsJson:         stateAsJson,
		},
	}
}
