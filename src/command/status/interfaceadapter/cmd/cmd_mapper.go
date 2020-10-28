package statuscmdadapter

import (
	"github.com/urfave/cli/v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/status"
	statuseventadapter "github.com/hekmekk/git-team/src/command/status/interfaceadapter/event"
	activation "github.com/hekmekk/git-team/src/shared/activation/impl"
	config "github.com/hekmekk/git-team/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
	state "github.com/hekmekk/git-team/src/shared/state/impl"
)

// Command the status command
func Command() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Print the current status",
		Action: func(c *cli.Context) error {
			return commandadapter.RunUrFave(Policy(), statuseventadapter.MapEventToEffects)(c)
		},
	}
}

// Policy the status policy constructor
func Policy() status.Policy {
	return status.Policy{
		Deps: status.Dependencies{
			ConfigReader:        config.NewGitconfigDataSource(gitconfig.NewDataSource()),
			StateReader:         state.NewGitConfigDataSource(gitconfig.NewDataSource()),
			ActivationValidator: activation.NewGitConfigDataSource(gitconfig.NewDataSource()),
		},
	}
}
