package removecmdadapter

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/assignments/remove"
	removeeventadapter "github.com/hekmekk/git-team/src/command/assignments/remove/interfaceadapter/event"
	gitconfiglegacy "github.com/hekmekk/git-team/src/shared/gitconfig/impl/legacy"
)

// Command the rm command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	rm := root.Command("rm", "Remove an alias to co-author assignment")
	alias := rm.Arg("alias", "The alias to identify the assignment to be removed").Required().String()

	rm.Action(commandadapter.Run(policy(alias), removeeventadapter.MapEventToEffects))

	return rm
}

func policy(alias *string) remove.Policy {
	return remove.Policy{
		Req: remove.DeAllocationRequest{
			Alias: alias,
		},
		Deps: remove.Dependencies{
			GitRemoveAlias: func(alias string) error {
				return gitconfiglegacy.UnsetAll(fmt.Sprintf("team.alias.%s", alias))
			},
		},
	}
}
