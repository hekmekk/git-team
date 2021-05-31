package main

import (
	"errors"
	"log"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/src/shared/cli/effects"

	addcmdadapter "github.com/hekmekk/git-team/src/command/assignments/add/cliadapter/cmd"
	assignmentscmdadapter "github.com/hekmekk/git-team/src/command/assignments/cliadapter/cmd"
	listcmdadapter "github.com/hekmekk/git-team/src/command/assignments/list/cliadapter/cmd"
	removecmdadapter "github.com/hekmekk/git-team/src/command/assignments/remove/cliadapter/cmd"
	completioncmdadapter "github.com/hekmekk/git-team/src/command/completion/cliadapter/cmd"
	configcmdadapter "github.com/hekmekk/git-team/src/command/config/cliadapter/cmd"
	disablecmdadapter "github.com/hekmekk/git-team/src/command/disable/cliadapter/cmd"
	enablecmdadapter "github.com/hekmekk/git-team/src/command/enable/cliadapter/cmd"
	statuscmdadapter "github.com/hekmekk/git-team/src/command/status/cliadapter/cmd"
)

const (
	version     = "1.7.0"
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
					return effects.NewExitErrMsg(errors.New("failed to generate man page")).Run()
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
