package status

import (
	"github.com/hekmekk/git-team/src/core/events"
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/shared/config/interface"
	state "github.com/hekmekk/git-team/src/shared/state/interface"
)

// Dependencies the dependencies of the status Policy module
type Dependencies struct {
	StateReader  state.Reader
	ConfigReader config.Reader
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply show the current status of git-team
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	_, cfgReadErr := deps.ConfigReader.Read()
	if cfgReadErr != nil {
		return StateRetrievalFailed{Reason: cfgReadErr}
	}

	state, stateRepositoryQueryErr := deps.StateReader.Query(activationscope.Global)
	if stateRepositoryQueryErr != nil {
		return StateRetrievalFailed{Reason: stateRepositoryQueryErr}
	}

	return StateRetrievalSucceeded{State: state}
}
