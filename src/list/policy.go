package list

import (
	"github.com/hekmekk/git-team/src/core/assignment"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/gitconfig/error"
)

// Dependencies the dependencies of the list Policy module
type Dependencies struct {
	GitGetAssignments func() (map[string]string, error)
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply show the available assignments
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	aliasCoauthorMap, err := deps.GitGetAssignments()
	if err != nil && err.Error() != gitconfigerror.SectionOrKeyIsInvalid {
		return RetrievalFailed{Reason: err}
	}

	assignments := []assignment.Assignment{}

	for alias, coauthor := range aliasCoauthorMap {
		assignments = append(assignments, assignment.Assignment{Alias: alias, Coauthor: coauthor})
	}

	return RetrievalSucceeded{Assignments: assignments}
}
