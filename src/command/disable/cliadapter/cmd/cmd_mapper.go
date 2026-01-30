package disablecmdadapter

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/v2/src/command/disable"
	disableeventadapter "github.com/hekmekk/git-team/v2/src/command/disable/cliadapter/event"
	statuscmdmapper "github.com/hekmekk/git-team/v2/src/command/status/cliadapter/cmd"
	activation "github.com/hekmekk/git-team/v2/src/shared/activation/impl"
	commandadapter "github.com/hekmekk/git-team/v2/src/shared/cli/commandadapter"
	configds "github.com/hekmekk/git-team/v2/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl"
	state "github.com/hekmekk/git-team/v2/src/shared/state/impl"
)

// Command the disable command
func Command() *cli.Command {
	return &cli.Command{
		Name:  "disable",
		Usage: "Use default commit template and remove prepare-commit-msg hook",
		Action: func(c *cli.Context) error {
			return commandadapter.Run(policy(), disableeventadapter.MapEventToEffectFactory(statuscmdmapper.Policy(false)))
		},
	}
}

func policy() disable.Policy {
	return disable.Policy{
		Deps: disable.Dependencies{
			ConfigReader:        configds.NewGitconfigDataSource(gitconfig.NewDataSource()),
			GitConfigReader:     gitconfig.NewDataSource(),
			GitConfigWriter:     gitconfig.NewDataSink(),
			StatFile:            os.Stat,
			RemoveFile:          os.RemoveAll,
			StateReader:         state.NewGitConfigDataSource(gitconfig.NewDataSource()),
			StateWriter:         state.NewGitConfigDataSink(gitconfig.NewDataSink()),
			ActivationValidator: activation.NewGitConfigDataSource(gitconfig.NewDataSource()),
		},
	}
}
