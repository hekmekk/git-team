package list

import (
	"github.com/hekmekk/git-team/src/core/assignment"
	"github.com/hekmekk/git-team/src/core/events"
)

// Dependencies the dependencies of the Policy
type Dependencies struct {
	GitGetAssignments func() (map[string]string, error)
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply assign a co-author to an alias
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	aliasCoauthorMap, err := deps.GitGetAssignments()
	if err != nil {
		return RetrievalFailed{Reason: err}
	}

	var assignments []assignment.Assignment

	for alias, coauthor := range aliasCoauthorMap {
		assignments = append(assignments, assignment.Assignment{Alias: alias, Coauthor: coauthor})
	}

	return RetrievalSucceeded{Assignments: assignments}
}
