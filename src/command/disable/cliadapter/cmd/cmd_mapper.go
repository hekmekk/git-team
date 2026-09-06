package disablecmdadapter

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
	"os"

	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/v2/src/command/disable"
	disableeventadapter "github.com/hekmekk/git-team/v2/src/command/disable/cliadapter/event"
	statuscmdmapper "github.com/hekmekk/git-team/v2/src/command/status/cliadapter/cmd"
	activation "github.com/hekmekk/git-team/v2/src/shared/activation/impl"
	commandadapter "github.com/hekmekk/git-team/v2/src/shared/cli/commandadapter"
	configds "github.com/hekmekk/git-team/v2/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl"
	state "github.com/hekmekk/git-team/v2/src/shared/state/impl"
)

// Command the disable command
func Command() *cli.Command {
	return &cli.Command{
		Name:  "disable",
		Usage: "Use default commit template and remove prepare-commit-msg hook",
		Action: func(c *cli.Context) error {
			return commandadapter.Run(policy(), disableeventadapter.MapEventToEffectFactory(statuscmdmapper.Policy(false)))
		},
	}
}

func policy() disable.Policy {
	return disable.Policy{
		Deps: disable.Dependencies{
			ConfigReader:        configds.NewGitconfigDataSource(gitconfig.NewDataSource()),
			GitConfigReader:     gitconfig.NewDataSource(),
			GitConfigWriter:     gitconfig.NewDataSink(),
			StatFile:            os.Stat,
			RemoveFile:          os.RemoveAll,
			StateReader:         state.NewGitConfigDataSource(gitconfig.NewDataSource()),
			StateWriter:         state.NewGitConfigDataSink(gitconfig.NewDataSink()),
			ActivationValidator: activation.NewGitConfigDataSource(gitconfig.NewDataSource()),
		},
	}
}
