package config

import (
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
		return RetrievalFailed{Reason: err}
	}

	return RetrievalSucceeded{Config: cfg}
}
