package add

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hekmekk/git-team/src/gitconfig"
	"github.com/hekmekk/git-team/src/validation"
)

// Definition the command, arguments, and dependencies
type Definition struct {
	CommandName string
	Args        Args
	Deps        Dependencies
}

// Args the arguments of the Runner
type Args struct {
	Alias    *string
	Coauthor *string
}

// Dependencies the dependencies of the Runner
type Dependencies struct {
	SanityCheckCoauthor func(string) error
	GitAddAlias         func(string, string) error
	GitResolveAlias     func(string) (string, error)
	StdinNewReader      func() *bufio.Reader
	StdinReadLine       func(reader *bufio.Reader) (string, error)
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
			SanityCheckCoauthor: validation.SanityCheckCoauthor,
			GitAddAlias:         gitconfig.AddAlias,
			GitResolveAlias:     gitconfig.ResolveAlias,
			StdinNewReader:      func() *bufio.Reader { return bufio.NewReader(os.Stdin) },
			StdinReadLine:       func(reader *bufio.Reader) (string, error) { return reader.ReadString('\n') },
		},
	}
}

const (
	y   string = "y"
	yes string = "yes"
)

// Apply assign a co-author to an alias
func Apply(deps Dependencies, args Args) interface{} {
	alias := *args.Alias
	coauthor := *args.Coauthor

	checkErr := deps.SanityCheckCoauthor(coauthor)
	if checkErr != nil {
		return AssignmentFailed{Reason: checkErr}
	}

	shouldAskForOverride, existingCoauthor := findExistingCoauthor(deps, alias)

	shouldAddAssignment := true
	if shouldAskForOverride {
		choice, err := shouldAssignmentBeOverridden(deps, alias, existingCoauthor, coauthor)
		if err != nil {
			return AssignmentFailed{Reason: err}
		}
		shouldAddAssignment = choice
	}

	switch shouldAddAssignment {
	case true:
		err := deps.GitAddAlias(alias, coauthor)
		if err != nil {
			return AssignmentFailed{Reason: err}
		}
		return AssignmentSucceeded{Alias: alias, Coauthor: coauthor}
	default:
		return AssignmentAborted{
			Alias:             alias,
			ExistingCoauthor:  existingCoauthor,
			ReplacingCoauthor: coauthor,
		}
	}
}

func shouldAssignmentBeOverridden(deps Dependencies, alias, existingCoauthor, replacingCoauthor string) (bool, error) {
	reader := deps.StdinNewReader()

	question := fmt.Sprintf("Alias '%s' -> '%s' exists already. Override with '%s'?", alias, existingCoauthor, replacingCoauthor)
	os.Stdout.WriteString(fmt.Sprintf("%s [N/y] ", question)) // ignoring errors for now, unlikely

	answer, readErr := deps.StdinReadLine(reader)
	if readErr != nil {
		return false, readErr
	}
	answer = strings.ToLower(strings.TrimSpace(strings.TrimRight(answer, "\n")))
	switch answer {
	case y, yes:
		return true, nil
	default:
		return false, nil
	}
}

func findExistingCoauthor(deps Dependencies, alias string) (bool, string) {
	existingCoauthor, resolveErr := deps.GitResolveAlias(alias)
	if resolveErr == nil {
		return true, existingCoauthor
	}

	return false, ""
}
