package disablecmdadapter

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/disable"
	disableeventadapter "github.com/hekmekk/git-team/src/command/disable/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/gitconfig"
	staterepository "github.com/hekmekk/git-team/src/core/state_repository"
	commitsettingsds "github.com/hekmekk/git-team/src/shared/commitsettings/datasource"
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
			GitUnsetCommitTemplate: func() error { return gitconfig.UnsetAll("commit.template") },
			GitUnsetHooksPath:      func() error { return gitconfig.UnsetAll("core.hooksPath") },
			CommitSettingsReader:   commitsettingsds.NewStaticValueDataSource(),
			StatFile:               os.Stat,
			RemoveFile:             os.Remove,
			PersistDisabled:        staterepository.PersistDisabled,
		},
	}
}
