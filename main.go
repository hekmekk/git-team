package main

import (
	"errors"
	"log"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/src/core/effects"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	addcmdadapter "github.com/hekmekk/git-team/src/command/assignments/add/interfaceadapter/cmd"
	assignmentscmdadapter "github.com/hekmekk/git-team/src/command/assignments/interfaceadapter/cmd"
	listcmdadapter "github.com/hekmekk/git-team/src/command/assignments/list/interfaceadapter/cmd"
	removecmdadapter "github.com/hekmekk/git-team/src/command/assignments/remove/interfaceadapter/cmd"
	completioncmdadapter "github.com/hekmekk/git-team/src/command/completion/interfaceadapter/cmd"
	configcmdadapter "github.com/hekmekk/git-team/src/command/config/interfaceadapter/cmd"
	disablecmdadapter "github.com/hekmekk/git-team/src/command/disable/interfaceadapter/cmd"
	enablecmdadapter "github.com/hekmekk/git-team/src/command/enable/interfaceadapter/cmd"
	statuscmdadapter "github.com/hekmekk/git-team/src/command/status/interfaceadapter/cmd"
)

const (
	version     = "1.6.0"
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
					return commandadapter.RunEffect(effects.NewExitErr(errors.New("failed to generate man page")))
				}
				return commandadapter.RunEffect(effects.NewExitOkMsg(manPage))
			}

			return enablecmdadapter.Command().Action(c)
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
