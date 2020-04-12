package configcmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	configds "github.com/hekmekk/git-team/src/command/config/datasource"
	configeventadapter "github.com/hekmekk/git-team/src/command/config/interfaceadapter/event"
	configpolicy "github.com/hekmekk/git-team/src/command/config/policy"
)

// Command the config command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	config := root.Command("config", "Edit configuration")

	config.Action(commandadapter.Run(policy(), configeventadapter.MapEventToEffects))

	return config
}

func policy() configpolicy.Policy {
	return configpolicy.Policy{
		Deps: configpolicy.Dependencies{
			ConfigReader: configds.NewGitconfigDataSource(),
		},
	}
}
