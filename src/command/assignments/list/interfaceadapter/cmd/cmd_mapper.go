package listcmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/assignments/list"
	listeventadapter "github.com/hekmekk/git-team/src/command/assignments/list/interfaceadapter/event"
	gitconfiglegacy "github.com/hekmekk/git-team/src/shared/gitconfig/impl/legacy"
	"github.com/hekmekk/git-team/src/shared/gitconfig/scope"
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
			GitGetAssignments: func() (map[string]string, error) { return gitconfiglegacy.GetRegexp(scope.Global, "team.alias") },
		},
	}
}
