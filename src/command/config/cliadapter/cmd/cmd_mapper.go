package configcmdadapter

import (
	"github.com/urfave/cli/v2"

	configeventadapter "github.com/hekmekk/git-team/src/command/config/cliadapter/event"
	configpolicy "github.com/hekmekk/git-team/src/command/config/policy"
	commandadapter "github.com/hekmekk/git-team/src/shared/cli/commandadapter"
	configdatasink "github.com/hekmekk/git-team/src/shared/config/datasink"
	configdatasource "github.com/hekmekk/git-team/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
)

// Command the config command
func Command() *cli.Command {
	return &cli.Command{
		Name:      "config",
		Usage:     "Edit configuration",
		ArgsUsage: "<key> <value>",
		Action: func(c *cli.Context) error {
			args := c.Args()
			key := args.First()
			value := args.Get(1)
			return commandadapter.Run(policy(&key, &value), configeventadapter.MapEventToEffect)
		},
	}
}

func policy(key *string, value *string) configpolicy.Policy {
	return configpolicy.Policy{
		Req: configpolicy.Request{
			Key:   key,
			Value: value,
		},
		Deps: configpolicy.Dependencies{
			ConfigReader: configdatasource.NewGitconfigDataSource(gitconfig.NewDataSource()),
			ConfigWriter: configdatasink.NewGitconfigDataSink(gitconfig.NewDataSink()),
		},
	}
}
