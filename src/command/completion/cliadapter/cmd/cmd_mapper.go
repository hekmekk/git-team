package completioncmdadapter

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

	"github.com/hekmekk/git-team/v2/src/command/completion/script"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
)

func Command() *cli.Command {
	bashCommand := &cli.Command{
		Name:        "bash",
		Usage:       "Bash completion",
		Description: "Source with bash to get auto completion",
		Action: func(c *cli.Context) error {
			return effects.NewExitOkMsg(script.Bash).Run()
		},
	}

	zshCommand := &cli.Command{
		Name:        "zsh",
		Usage:       "Zsh completion",
		Description: "Source with zsh to get auto completion",
		Action: func(c *cli.Context) error {
			return effects.NewExitOkMsg(script.Zsh).Run()
		},
	}

	return &cli.Command{
		Name:        "completion",
		Usage:       "Shell completion",
		Description: "Source the output of this command to get auto completion",
		Subcommands: []*cli.Command{
			bashCommand,
			zshCommand,
		},
	}
}
