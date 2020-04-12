package policy

import (
	configevents "github.com/hekmekk/git-team/src/command/config/events"
	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/events"
)

// Dependencies the dependencies of the config Policy module
type Dependencies struct {
	ConfigRepository config.Repository
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply show the current status of git-team
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	cfg, err := deps.ConfigRepository.Query()
	if err != nil {
		return configevents.RetrievalFailed{Reason: err}
	}

	return configevents.RetrievalSucceeded{Config: cfg}
}
