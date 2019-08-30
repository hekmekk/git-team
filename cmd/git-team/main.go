package main

import (
	"os"

	"github.com/hekmekk/git-team/src/add/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/add/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/policy"
	"github.com/hekmekk/git-team/src/disable/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/disable/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/enable/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/enable/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/list/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/list/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/remove/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/remove/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/status/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/status/interfaceadapter/event"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "v1.3.2"
	author  = "Rea Sand <hekmek@posteo.de>"
)

func main() {
	application := newApplication(author, version)

	switch kingpin.MustParse(application.app.Parse(os.Args[1:])) {
	case application.add.CommandName:
		applyPolicy(application.add.Policy, addeventadapter.MapEventToEffects)
	case application.remove.CommandName:
		applyPolicy(application.remove.Policy, removeeventadapter.MapEventToEffects)
	case application.enable.CommandName:
		applyPolicy(application.enable.Policy, enableeventadapter.MapEventToEffectsFactory(application.status.Policy.Deps.StateRepositoryQuery))
	case application.disable.CommandName:
		applyPolicy(application.disable.Policy, disableeventadapter.MapEventToEffectsFactory(application.status.Policy.Deps.StateRepositoryQuery))
	case application.status.CommandName:
		applyPolicy(application.status.Policy, statuseventadapter.MapEventToEffects)
	case application.list.CommandName:
		applyPolicy(application.list.Policy, listeventadapter.MapEventToEffects)
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
	app     *kingpin.Application
	add     addcmdadapter.Definition
	remove  removecmdadapter.Definition
	enable  enablecmdadapter.Definition
	disable disablecmdadapter.Definition
	status  statuscmdadapter.Definition
	list    listcmdadapter.Definition
}

func newApplication(author string, version string) application {
	app := kingpin.New("git-team", "Command line interface for managing and enhancing git commit messages with co-authors.")

	app.HelpFlag.Short('h')
	app.Version(version)
	app.Author(author)

	return application{
		app:     app,
		add:     addcmdadapter.NewDefinition(app),
		remove:  removecmdadapter.NewDefinition(app),
		enable:  enablecmdadapter.NewDefinition(app),
		disable: disablecmdadapter.NewDefinition(app),
		status:  statuscmdadapter.NewDefinition(app),
		list:    listcmdadapter.NewDefinition(app),
	}
}
