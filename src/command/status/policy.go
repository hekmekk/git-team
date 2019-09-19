package status

import (
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/state"
)

// Dependencies the dependencies of the status Policy module
type Dependencies struct {
	StateRepositoryQuery func() (state.State, error)
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply show the current status of git-team
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	state, err := deps.StateRepositoryQuery()
	if err != nil {
		return StateRetrievalFailed{Reason: err}
	}

	return StateRetrievalSucceeded{State: state}
}
