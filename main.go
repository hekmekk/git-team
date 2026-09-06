package main

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
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"

	addcmdadapter "github.com/hekmekk/git-team/v2/src/command/assignments/add/cliadapter/cmd"
	assignmentscmdadapter "github.com/hekmekk/git-team/v2/src/command/assignments/cliadapter/cmd"
	listcmdadapter "github.com/hekmekk/git-team/v2/src/command/assignments/list/cliadapter/cmd"
	removecmdadapter "github.com/hekmekk/git-team/v2/src/command/assignments/remove/cliadapter/cmd"
	completioncmdadapter "github.com/hekmekk/git-team/v2/src/command/completion/cliadapter/cmd"
	configcmdadapter "github.com/hekmekk/git-team/v2/src/command/config/cliadapter/cmd"
	disablecmdadapter "github.com/hekmekk/git-team/v2/src/command/disable/cliadapter/cmd"
	enablecmdadapter "github.com/hekmekk/git-team/v2/src/command/enable/cliadapter/cmd"
	statuscmdadapter "github.com/hekmekk/git-team/v2/src/command/status/cliadapter/cmd"
)

const (
	version     = "2.0.0"
	authorName  = "Rea Sand"
	authorEmail = "hekmek@posteo.de"
)

func main() {
	application := newApplication()
	err := application.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func newApplication() *cli.App {
	app := &cli.App{
		Name:     "git-team",
		Compiled: time.Now(),
		Version:  version,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  authorName,
				Email: authorEmail,
			},
		},
		Usage:                "Command line interface for managing and enhancing git commit messages with co-authors.",
		EnableBashCompletion: true,
		HideHelp:             false,
		HideVersion:          false,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "generate-man-page", Value: false, Usage: "Generate man page for this"},
		},
		Commands: []*cli.Command{
			enablecmdadapter.Command(),
			disablecmdadapter.Command(),
			statuscmdadapter.Command(),
			assignmentscmdadapter.Command(),
			addcmdadapter.Command(),
			listcmdadapter.Command(),
			removecmdadapter.Command(),
			configcmdadapter.Command(),
			completioncmdadapter.Command(),
		},
		Action: func(c *cli.Context) error {
			shouldGenerateManPage := c.Bool("generate-man-page")
			if shouldGenerateManPage {
				manPage, err := c.App.ToMan()
				if err != nil {
					return effects.NewExitErrMsg(fmt.Errorf("failed to generate man page: %s", err)).Run()
				}
				return effects.NewExitOkMsg(manPage).Run()
			}

			return enablecmdadapter.Command().Action(c)
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
