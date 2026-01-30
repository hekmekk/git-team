package statuscmdadapter

import (
	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/v2/src/command/status"
	statuseventadapter "github.com/hekmekk/git-team/v2/src/command/status/cliadapter/event"
	activation "github.com/hekmekk/git-team/v2/src/shared/activation/impl"
	commandadapter "github.com/hekmekk/git-team/v2/src/shared/cli/commandadapter"
	config "github.com/hekmekk/git-team/v2/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl"
	state "github.com/hekmekk/git-team/v2/src/shared/state/impl"
)

// Command the status command
func Command() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Print the current status",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "json", Value: false, Usage: "Display the current status as a JSON formatted struct"},
		},
		Action: func(c *cli.Context) error {
			stateAsJson := c.Bool("json")

			return commandadapter.Run(Policy(stateAsJson), statuseventadapter.MapEventToEffect)
		},
	}
}

// Policy the status policy constructor
func Policy(stateAsJson bool) status.Policy {
	return status.Policy{
		Deps: status.Dependencies{
			ConfigReader:        config.NewGitconfigDataSource(gitconfig.NewDataSource()),
			StateReader:         state.NewGitConfigDataSource(gitconfig.NewDataSource()),
			ActivationValidator: activation.NewGitConfigDataSource(gitconfig.NewDataSource()),
			StateAsJson:         stateAsJson,
		},
	}
}
