package configcmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	configdatasink "github.com/hekmekk/git-team/src/command/config/datasink"
	configdatasource "github.com/hekmekk/git-team/src/command/config/datasource"
	configeventadapter "github.com/hekmekk/git-team/src/command/config/interfaceadapter/event"
	configpolicy "github.com/hekmekk/git-team/src/command/config/policy"
)

// Command the config command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	config := root.Command("config", "Edit configuration")
	key := config.Arg("key", "the key of the setting").String()
	value := config.Arg("value", "the value of the setting").String()

	config.Action(commandadapter.Run(policy(key, value), configeventadapter.MapEventToEffects))

	return config
}

func policy(key *string, value *string) configpolicy.Policy {
	return configpolicy.Policy{
		Req: configpolicy.Request{
			Key:   key,
			Value: value,
		},
		Deps: configpolicy.Dependencies{
			ConfigReader: configdatasource.NewGitconfigDataSource(),
			ConfigWriter: configdatasink.NewGitconfigDataSink(),
		},
	}
}
