package add

import (
	"fmt"

	"github.com/hekmekk/git-team/src/effects"
	"github.com/hekmekk/git-team/src/gitconfig"
	"github.com/hekmekk/git-team/src/validation"

	"github.com/fatih/color"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Args the arguments of the Runner
type Args struct {
	Alias    *string
	Coauthor *string
}

// Dependencies the dependencies of the Runner
type Dependencies struct {
	AddGitAlias func(string, string) error
}

// Definition the command, arguments, and dependencies
type Definition struct {
	Command *kingpin.CmdClause
	Args    Args
	Deps    Dependencies
}

// New the constructor for Definition
func New(app *kingpin.Application) Definition {
	command := app.Command("add", "Add an alias")
	return Definition{
		Command: command,
		Args: Args{
			Alias:    command.Arg("alias", "The alias to be added").Required().String(),
			Coauthor: command.Arg("coauthor", "The co-author").Required().String(),
		},
		Deps: Dependencies{
			AddGitAlias: gitconfig.AddAlias,
		},
	}
}

// Run assign a co-author to an alias
func Run(deps Dependencies, args Args) []effects.Effect {
	err := assignCoauthorToAlias(deps, args)
	if err != nil {
		return []effects.Effect{
			effects.NewPrintErr(err),
			effects.NewExitErr(),
		}
	}

	return []effects.Effect{
		effects.NewPrintMessage(color.GreenString(fmt.Sprintf("Alias '%s' -> '%s' has been added.", *args.Alias, *args.Coauthor))),
		effects.NewExitOk(),
	}
}

func assignCoauthorToAlias(deps Dependencies, args Args) error {
	checkErr := validation.SanityCheckCoauthor(*args.Coauthor) // TODO: Add as dependency. Has no side effects but doesn't need to be tested here.
	if checkErr != nil {
		return checkErr
	}

	addErr := deps.AddGitAlias(*args.Alias, *args.Coauthor)
	if addErr != nil {
		return addErr
	}
	return nil
}
