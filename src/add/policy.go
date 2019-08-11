package add

import (
	"fmt"
	"os" // TODO: remove this dependency aswell. dep: question: string -> answer: string
	"strings"
)

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
