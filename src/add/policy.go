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
	Request     AssignmentRequest
	Deps        Dependencies
}

// AssignmentRequest the arguments of the Runner
type AssignmentRequest struct {
	Alias    *string
	Coauthor *string
}

// Dependencies the dependencies of the Runner
type Dependencies struct {
	SanityCheckCoauthor func(string) error
	GitAddAlias         func(string, string) error
	GitResolveAlias     func(string) (string, error)
	StdinReadLine       func() (string, error)
}

// New the constructor for Definition
func New(name string, alias, coauthor *string) Definition {
	return Definition{
		CommandName: name,
		Request: AssignmentRequest{
			Alias:    alias,
			Coauthor: coauthor,
		},
		Deps: Dependencies{
			SanityCheckCoauthor: validation.SanityCheckCoauthor,
			GitAddAlias:         gitconfig.AddAlias,
			GitResolveAlias:     gitconfig.ResolveAlias,
			StdinReadLine:       func() (string, error) { return bufio.NewReader(os.Stdin).ReadString('\n') },
		},
	}
}

const (
	y   string = "y"
	yes string = "yes"
)

// Apply assign a co-author to an alias
func Apply(deps Dependencies, req AssignmentRequest) interface{} {
	alias := *req.Alias
	coauthor := *req.Coauthor

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
	question := fmt.Sprintf("Alias '%s' -> '%s' exists already. Override with '%s'?", alias, existingCoauthor, replacingCoauthor)
	os.Stdout.WriteString(fmt.Sprintf("%s [N/y] ", question)) // ignoring errors for now, unlikely

	answer, readErr := deps.StdinReadLine()
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
