package configcmdadapter

import (
	"fmt"

	"github.com/urfave/cli/v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	configeventadapter "github.com/hekmekk/git-team/src/command/config/interfaceadapter/event"
	configpolicy "github.com/hekmekk/git-team/src/command/config/policy"
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
			return commandadapter.RunUrFave(policy(&key, &value), configeventadapter.MapEventToEffects)(c)
		},
		BashComplete: func(c *cli.Context) {
			options := map[string][]string{
				"activation-scope": []string{"repo-local", "global"},
			}

			args := c.Args()
			argsLen := args.Len()

			if argsLen == 0 {
				for key := range options {
					fmt.Println(key)
				}
			}

			if argsLen == 1 {
				values := options[args.First()]
				for _, value := range values {
					fmt.Println(value)
				}
			}
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
