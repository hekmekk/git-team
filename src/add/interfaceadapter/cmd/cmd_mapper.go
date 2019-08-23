package addcmdadapter

import (
	"bufio"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/add"
	"github.com/hekmekk/git-team/src/gitconfig"
	"github.com/hekmekk/git-team/src/validation"
)

// Definition the command, arguments, and dependencies
type Definition struct {
	CommandName string
	Request     add.AssignmentRequest
	Deps        add.Dependencies
}

// NewDefinition the constructor for Definition
func NewDefinition(app *kingpin.Application) Definition {

	command := app.Command("add", "Add a new or override an existing alias to co-author assignment")
	alias := command.Arg("alias", "The alias to assign a co-author to").Required().String()
	coauthor := command.Arg("coauthor", "The co-author").Required().String()

	// TODO: create NewPolicy here?!
	return Definition{
		CommandName: command.FullCommand(),
		Request: add.AssignmentRequest{
			Alias:    alias,
			Coauthor: coauthor,
		},
		Deps: add.Dependencies{
			SanityCheckCoauthor: validation.SanityCheckCoauthor,
			GitResolveAlias:     gitconfig.ResolveAlias,
			GitAddAlias:         gitconfig.AddAlias,
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
