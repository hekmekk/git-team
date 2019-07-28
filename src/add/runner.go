package add

import (
	"fmt"

	"github.com/hekmekk/git-team/src/effects"
	"github.com/hekmekk/git-team/src/gitconfig"

	"github.com/fatih/color"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Definition the command and its arguments
type Definition struct {
	Command *kingpin.CmdClause
	Args    Args
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
	}
}

// Run run the add functionality
func Run(args Args) []effects.Effect {
	addDeps := Dependencies{
		AddGitAlias: gitconfig.AddAlias,
	}
	execAdd := ExecutorFactory(addDeps)

	err := execAdd(args)
	if err != nil {
		return []effects.Effect{effects.NewPrintErr(err), effects.NewExitErr()}
	}

	return []effects.Effect{effects.NewPrintMessage(color.GreenString(fmt.Sprintf("Alias '%s' -> '%s' has been added.", *args.Alias, *args.Coauthor))), effects.NewExitOk()}
}
