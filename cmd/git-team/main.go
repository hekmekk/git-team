package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/command/assignments/add/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/assignments/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/assignments/list/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/assignments/remove/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/disable/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/enable/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/status/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/core/effects"
)

const (
	version = "v1.3.6-alpha1"
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

	assignmentscmdadapter.Command(app)
	enablecmdadapter.Command(app)
	disablecmdadapter.Command(app)
	statuscmdadapter.Command(app)

	return app
}
