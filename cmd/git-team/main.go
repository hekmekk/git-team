package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/src/core/effects"

	addcmdadapter "github.com/hekmekk/git-team/src/command/assignments/add/interfaceadapter/cmd"
	assignmentscmdadapter "github.com/hekmekk/git-team/src/command/assignments/interfaceadapter/cmd"
	listcmdadapter "github.com/hekmekk/git-team/src/command/assignments/list/interfaceadapter/cmd"

	removecmdadapter "github.com/hekmekk/git-team/src/command/assignments/remove/interfaceadapter/cmd"
	configcmdadapter "github.com/hekmekk/git-team/src/command/config/interfaceadapter/cmd"
	disablecmdadapter "github.com/hekmekk/git-team/src/command/disable/interfaceadapter/cmd"
	enablecmdadapter "github.com/hekmekk/git-team/src/command/enable/interfaceadapter/cmd"
	statuscmdadapter "github.com/hekmekk/git-team/src/command/status/interfaceadapter/cmd"
)

const (
	version     = "1.5.2"
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
	ls := listcmdadapter.Command()
	ls.Before = func(c *cli.Context) error {
		effects.NewDeprecationWarning("git team ls", "git team assignments").Run()
		return nil
	}

	add := addcmdadapter.Command()
	add.Before = func(c *cli.Context) error {
		effects.NewDeprecationWarning("git team add", "git team assignments add").Run()
		return nil
	}

	rm := removecmdadapter.Command()
	rm.Before = func(c *cli.Context) error {
		effects.NewDeprecationWarning("git team rm", "git team assignments rm").Run()
		return nil
	}

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
			add,
			ls,
			rm,
			configcmdadapter.Command(),
		},
		Action: func(c *cli.Context) error {
			shouldGenerateManPage := c.Bool("generate-man-page")
			if shouldGenerateManPage {
				manPage, err := c.App.ToMan()
				if err != nil {
					effects.NewPrintErr(errors.New("failed to generate man page")).Run()
					return nil
				}
				fmt.Println(manPage)
				return nil
			}

			effects.NewDeprecationWarning("git team (as a default command)", "git team enable").Run()
			return enablecmdadapter.Command().Action(c)
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
