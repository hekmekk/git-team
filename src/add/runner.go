package add

import (
	"fmt"
	"github.com/hekmekk/git-team/src/gitconfig"
	"os"

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

type Effect interface {
	Run()
}

type PrintMessage struct {
	message string
}

func (printMsg PrintMessage) Run() {
	color.Green(printMsg.message)
}

type PrintErr struct {
	err error
}

func (printErr PrintErr) Run() {
	os.Stderr.WriteString(fmt.Sprintf("error: %s\n", printErr.err))
}

type ExitErr struct{}

func (_ ExitErr) Run() {
	os.Exit(-1)
}

// Run run the add functionality
func Run(args Args) []Effect {
	addDeps := Dependencies{
		AddGitAlias: gitconfig.AddAlias,
	}
	execAdd := ExecutorFactory(addDeps)

	err := execAdd(args)
	if err != nil {
		return []Effect{PrintErr{err: err}, ExitErr{}}
	}

	return []Effect{PrintMessage{message: fmt.Sprintf("Alias '%s' -> '%s' has been added.", *args.Alias, *args.Coauthor)}}
}
