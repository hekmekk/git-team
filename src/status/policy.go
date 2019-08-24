package status

import (
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/state"
)

// Dependencies the dependencies of the Policy
type Dependencies struct {
	StateRepositoryQuery func() (state.State, error)
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply assign a co-author to an alias
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	state, err := deps.StateRepositoryQuery()
	if err != nil {
		return StateRetrievalFailed{Reason: err}
	}

	return StateRetrievalSucceeded{State: state}
}
