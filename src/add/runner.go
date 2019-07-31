package add

import (
	"fmt"

	"github.com/hekmekk/git-team/src/effects"
	"github.com/hekmekk/git-team/src/gitconfig"
	"github.com/hekmekk/git-team/src/validation"

	"github.com/fatih/color"
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
	CommandName string
	Args        Args
	Deps        Dependencies
}

// New the constructor for Definition
func New(name string, alias, coauthor *string) Definition {
	return Definition{
		CommandName: name,
		Args: Args{
			Alias:    alias,
			Coauthor: coauthor,
		},
		Deps: Dependencies{
			AddGitAlias: gitconfig.AddAlias,
		},
	}
}

// Run assign a co-author to an alias
// TODO: add should not know stuff like PrintMessage or ExitOk -> return Events which get mapped to Effects?!
func Run(deps Dependencies, args Args) []effects.Effect {
	alias := *args.Alias
	coauthor := *args.Coauthor
	err := assignCoauthorToAlias(deps, alias, coauthor)
	if err != nil {
		return []effects.Effect{
			effects.NewPrintErr(err),
			effects.NewExitErr(),
		}
	}

	return []effects.Effect{
		effects.NewPrintMessage(color.GreenString(fmt.Sprintf("Alias '%s' -> '%s' has been added.", alias, coauthor))),
		effects.NewExitOk(),
	}
}

func assignCoauthorToAlias(deps Dependencies, alias string, coauthor string) error {
	checkErr := validation.SanityCheckCoauthor(coauthor) // TODO: Add as dependency. Has no side effects but doesn't need to be tested here.
	if checkErr != nil {
		return checkErr
	}

	addErr := deps.AddGitAlias(alias, coauthor)
	if addErr != nil {
		return addErr
	}
	return nil
}
