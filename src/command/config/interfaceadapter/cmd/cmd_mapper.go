package configcmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	configeventadapter "github.com/hekmekk/git-team/src/command/config/interfaceadapter/event"
	configpolicy "github.com/hekmekk/git-team/src/command/config/policy"
	configds "github.com/hekmekk/git-team/src/core/config"
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
			ConfigRepository: configds.NewStaticValueDataSource(),
		},
	}
}
