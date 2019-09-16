package listcmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/core/gitconfig"
	"github.com/hekmekk/git-team/src/list"
	"github.com/hekmekk/git-team/src/list/interfaceadapter/event"
)

// Command the ls command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	ls := root.Command("ls", "List your assignments")
	ls.Alias("list")

	ls.Action(commandadapter.Run(policy(), listeventadapter.MapEventToEffects))

	return ls
}

func policy() list.Policy {
	return list.Policy{
		Deps: list.Dependencies{
			GitGetAssignments: func() (map[string]string, error) { return gitconfig.GetRegexp("team.alias") },
		},
	}
}
