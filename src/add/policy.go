package add

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hekmekk/git-team/src/gitconfig"
	"github.com/hekmekk/git-team/src/validation"
)

// Args the arguments of the Runner
type Args struct {
	Alias    *string
	Coauthor *string
}

// Dependencies the dependencies of the Runner
type Dependencies struct {
	AddGitAlias     func(string, string) error
	GitResolveAlias func(string) (string, error)
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
			AddGitAlias:     gitconfig.AddAlias,
			GitResolveAlias: gitconfig.ResolveAlias,
		},
	}
}

// AssignmentFailed assignment of coauthor to alias failed with Reason
type AssignmentFailed struct {
	Reason error
}

// AssignmentSucceeded assignment of coauthor to alias succeeded
type AssignmentSucceeded struct {
	Alias    string
	Coauthor string
}

// AssignmentAborted nothing happened
type AssignmentAborted struct {
	Alias             string
	ExistingCoauthor  string
	ReplacingCoauthor string
}

const (
	y    string = "y"
	yes  string = "yes"
	n    string = "n"
	no   string = "no"
	none string = ""
)

// TODO: need to inject io handling for unit tests
// Apply assign a co-author to an alias
func Apply(deps Dependencies, args Args) interface{} {
	alias := *args.Alias
	coauthor := *args.Coauthor
	existingCoauthor := findExistingCoauthor(deps, alias)

	if "" != existingCoauthor {
		reader := bufio.NewReader(os.Stdin)
		question := fmt.Sprintf("Alias '%s' -> '%s' exists already. Override with '%s'? [N/y] ", alias, existingCoauthor, coauthor)
		fmt.Print(question)
		answer, readErr := reader.ReadString('\n')
		if readErr != nil {
			return AssignmentFailed{Reason: readErr}
		}
		answer = strings.ToLower(strings.TrimSpace(strings.TrimRight(answer, "\n")))
		switch answer {
		case y, yes:
			err := assignCoauthorToAlias(deps, alias, coauthor)
			if err != nil {
				return AssignmentFailed{Reason: err}
			}
			return AssignmentSucceeded{Alias: alias, Coauthor: coauthor}
		case n, no, none:
			return AssignmentAborted{
				Alias:             alias,
				ExistingCoauthor:  existingCoauthor,
				ReplacingCoauthor: coauthor,
			}
		default:
			return AssignmentFailed{Reason: errors.New("invalid answer :|")}
		}
	}
	err := assignCoauthorToAlias(deps, alias, coauthor)
	if err != nil {
		return AssignmentFailed{Reason: err}
	}

	return AssignmentSucceeded{Alias: alias, Coauthor: coauthor}
}

func shouldOverride(deps Dependencies, alias, existingCoauthor, replacingCoauthor string) bool {
	return false
}

func findExistingCoauthor(deps Dependencies, alias string) string {
	existingCoauthor, resolveErr := deps.GitResolveAlias(alias)
	if resolveErr == nil {
		return existingCoauthor
	}

	return ""
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
