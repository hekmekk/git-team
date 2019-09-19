package main

import (
	"os"

	"github.com/hekmekk/git-team/src/command/assignments/add/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/assignments/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/assignments/list/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/assignments/remove/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/disable/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/enable/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/status/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/command/status/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/policy"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "v1.3.5-alpha9"
	author  = "Rea Sand <hekmek@posteo.de>"
)

func main() {
	application := newApplication(author, version)

	switch kingpin.MustParse(application.app.Parse(os.Args[1:])) {
	case application.status.CommandName:
		applyPolicy(application.status.Policy, statuseventadapter.MapEventToEffects)
	}

	os.Exit(0)
}

func applyPolicy(policy policy.Policy, adapter func(events.Event) []effects.Effect) {
	effects := adapter(policy.Apply())
	for _, effect := range effects {
		effect.Run()
	}
}

type application struct {
	app    *kingpin.Application
	status statuscmdadapter.Definition
}

func newApplication(author string, version string) application {
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

	return application{
		app:    app, // TODO: use actions and just return this ...
		status: statuscmdadapter.NewDefinition(app),
	}
}
