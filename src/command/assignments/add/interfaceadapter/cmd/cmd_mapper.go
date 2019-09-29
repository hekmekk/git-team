package addcmdadapter

import (
	"bufio"
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/assignments/add"
	"github.com/hekmekk/git-team/src/command/assignments/add/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/gitconfig"
	"github.com/hekmekk/git-team/src/core/validation"
)

// Command the add command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	add := root.Command("add", "Add a new or override an existing alias to co-author assignment")
	alias := add.Arg("alias", "The alias to assign a co-author to").Required().String()
	coauthor := add.Arg("coauthor", "The co-author").Required().String()

	add.Action(commandadapter.Run(policy(alias, coauthor), addeventadapter.MapEventToEffects))

	return add
}

func policy(alias *string, coauthor *string) add.Policy {
	return add.Policy{
		Req: add.AssignmentRequest{
			Alias:    alias,
			Coauthor: coauthor,
		},
		Deps: add.Dependencies{
			SanityCheckCoauthor: validation.SanityCheckCoauthor,
			GitResolveAlias:     commandadapter.ResolveAlias,
			GitAddAlias: func(alias string, coauthor string) error {
				return gitconfig.ReplaceAll(fmt.Sprintf("team.alias.%s", alias), coauthor) // TODO: move this logic into the policy
			},
			GetAnswerFromUser: func(question string) (string, error) {
				_, err := os.Stdout.WriteString(question)
				if err != nil {
					return "", err
				}
				return bufio.NewReader(os.Stdin).ReadString('\n')
			},
		},
	}
}
