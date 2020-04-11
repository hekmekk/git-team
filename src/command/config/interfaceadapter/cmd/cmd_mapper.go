package configcmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/config"
	configeventadapter "github.com/hekmekk/git-team/src/command/config/interfaceadapter/event"
	configds "github.com/hekmekk/git-team/src/core/config"
)

// Command the config command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	config := root.Command("config", "Edit configuration")

	config.Action(commandadapter.Run(policy(), configeventadapter.MapEventToEffects))

	return config
}

func policy() config.Policy {
	return config.Policy{
		Deps: config.Dependencies{
			ConfigRepository: configds.NewStaticValueDataSource(),
		},
	}
}
