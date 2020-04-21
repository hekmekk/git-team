package statuscmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/status"
	statuseventadapter "github.com/hekmekk/git-team/src/command/status/interfaceadapter/event"
	config "github.com/hekmekk/git-team/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
	state "github.com/hekmekk/git-team/src/shared/state/impl"
)

// Command the status command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	status := root.Command("status", "Print the current status")

	status.Action(commandadapter.Run(Policy(), statuseventadapter.MapEventToEffects))

	return status
}

// Policy the status policy constructor
func Policy() status.Policy {
	return status.Policy{
		Deps: status.Dependencies{
			ConfigReader: config.NewGitconfigDataSource(gitconfig.NewDataSource()),
			StateReader:  state.NewGitConfigDataSource(gitconfig.NewDataSource()),
		},
	}
}
