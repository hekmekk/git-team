package listcmdadapter

import (
	"gopkg.in/alecthomas/kingpin.v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/assignments/list"
	listeventadapter "github.com/hekmekk/git-team/src/command/assignments/list/interfaceadapter/event"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
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
			GitConfigReader: gitconfig.NewDataSource(),
		},
	}
}
