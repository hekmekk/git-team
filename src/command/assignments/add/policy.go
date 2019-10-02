package add

import (
	"fmt"
	"strings"

	"github.com/hekmekk/git-team/src/core/events"
)

// AssignmentRequest which coauthor to assign to the alias
type AssignmentRequest struct {
	Alias         *string
	Coauthor      *string
	ForceOverride *bool
}

// Dependencies the dependencies of the add Policy module
type Dependencies struct {
	SanityCheckCoauthor func(string) error
	GitAddAlias         func(string, string) error
	GitResolveAlias     func(string) (string, error)
	GetAnswerFromUser   func(string) (string, error)
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
	Req  AssignmentRequest
}

const (
	y   string = "y"
	yes string = "yes"
)

// Apply assign a co-author to an alias
func (policy Policy) Apply() events.Event {
	req := policy.Req
	deps := policy.Deps

	alias := *req.Alias
	coauthor := *req.Coauthor

	if checkErr := deps.SanityCheckCoauthor(coauthor); checkErr != nil {
		return AssignmentFailed{Reason: checkErr}
	}

	isAssignmentExisting := false

	existingCoauthor, resolveErr := deps.GitResolveAlias(alias)
	if resolveErr == nil {
		isAssignmentExisting = true
	} else {
		isAssignmentExisting = false
	}

	shouldAddAssignment := true

	if isAssignmentExisting && !*req.ForceOverride {
		choice, err := shouldAssignmentBeOverridden(deps, alias, existingCoauthor, coauthor)
		if err != nil {
			return AssignmentFailed{Reason: err}
		}
		shouldAddAssignment = choice
	}

	if !shouldAddAssignment {
		return AssignmentAborted{
			Alias:             alias,
			ExistingCoauthor:  existingCoauthor,
			ReplacingCoauthor: coauthor,
		}
	}

	err := deps.GitAddAlias(alias, coauthor)
	if err != nil {
		return AssignmentFailed{Reason: err}
	}

	return AssignmentSucceeded{Alias: alias, Coauthor: coauthor}
}

func shouldAssignmentBeOverridden(deps Dependencies, alias, existingCoauthor, replacingCoauthor string) (bool, error) {
	question := fmt.Sprintf("Assignment '%s' →  '%s' exists already. Override with '%s'? [N/y] ", alias, existingCoauthor, replacingCoauthor)

	answer, err := deps.GetAnswerFromUser(question)
	if err != nil {
		return false, err
	}

	answer = strings.ToLower(strings.TrimSpace(strings.TrimRight(answer, "\n")))
	switch answer {
	case y, yes:
		return true, nil
	default:
		return false, nil
	}
}
