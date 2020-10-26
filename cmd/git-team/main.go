package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	addcmdadapter "github.com/hekmekk/git-team/src/command/assignments/add/interfaceadapter/cmd"
	assignmentscmdadapter "github.com/hekmekk/git-team/src/command/assignments/interfaceadapter/cmd"
	listcmdadapter "github.com/hekmekk/git-team/src/command/assignments/list/interfaceadapter/cmd"
	removecmdadapter "github.com/hekmekk/git-team/src/command/assignments/remove/interfaceadapter/cmd"
	configcmdadapter "github.com/hekmekk/git-team/src/command/config/interfaceadapter/cmd"
	disablecmdadapter "github.com/hekmekk/git-team/src/command/disable/interfaceadapter/cmd"
	enablecmdadapter "github.com/hekmekk/git-team/src/command/enable/interfaceadapter/cmd"
	statuscmdadapter "github.com/hekmekk/git-team/src/command/status/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/core/effects"
)

const (
	version = "v1.4.2-rc1"
	author  = "Rea Sand <hekmek@posteo.de>"
)

func main() {
	application := newApplication(author, version)
	kingpin.MustParse(application.Parse(os.Args[1:]))
	os.Exit(0)
}

func newApplication(author string, version string) *kingpin.Application {
	app := kingpin.New("git-team", "Command line interface for managing and enhancing git commit messages with co-authors.")

	app.Author(author)
	app.Version(version)

	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	ls := listcmdadapter.Command(app)
	ls.PreAction(func(c *kingpin.ParseContext) error {
		effects.NewDeprecationWarning("git team ls", "git team assignments").Run()
		return nil
	})

	add := addcmdadapter.Command(app)
	add.PreAction(func(c *kingpin.ParseContext) error {
		effects.NewDeprecationWarning("git team add", "git team assignments add").Run()
		return nil
	})

	rm := removecmdadapter.Command(app)
	rm.PreAction(func(c *kingpin.ParseContext) error {
		effects.NewDeprecationWarning("git team rm", "git team assignments rm").Run()
		return nil
	})

	enable := enablecmdadapter.Command(app).Default()
	enable.PreAction(func(c *kingpin.ParseContext) error {
		index := c.Peek().Index
		numElements := len(c.Elements)
		if index == 1 && numElements == 1 {
			effects.NewDeprecationWarning("git team enable (without aliases)", "git team [status]").Run()
		}
		if index >= 1 && numElements == index+1 {
			effects.NewDeprecationWarning("git team (without further sub-command specification)", "git team enable").Run()
		}
		return nil
	})

	assignmentscmdadapter.Command(app)
	disablecmdadapter.Command(app)
	statuscmdadapter.Command(app)
	configcmdadapter.Command(app)

	return app
}
