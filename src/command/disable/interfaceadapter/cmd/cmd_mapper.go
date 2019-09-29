package disablecmdadapter

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/disable"
	"github.com/hekmekk/git-team/src/command/disable/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/gitconfig"
	"github.com/hekmekk/git-team/src/core/state_repository"
)

// Command the disable command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	disable := root.Command("disable", "Use default commit template and remove prepare-commit-msg hook")

	disable.Action(commandadapter.Run(policy(), disableeventadapter.MapEventToEffectsFactory(staterepository.Query)))

	return disable
}

func policy() disable.Policy {
	return disable.Policy{
		Deps: disable.Dependencies{
			GitUnsetCommitTemplate: gitconfig.UnsetCommitTemplate,
			GitUnsetHooksPath:      func() error { return gitconfig.UnsetAll("core.hooksPath") }, // TODO: move this logic into the policy
			LoadConfig:             config.Load,
			StatFile:               os.Stat,
			RemoveFile:             os.Remove,
			PersistDisabled:        staterepository.PersistDisabled,
		},
	}
}
