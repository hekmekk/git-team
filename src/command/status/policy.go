package status

import (
	"fmt"

	"github.com/hekmekk/git-team/v2/src/core/events"
	activation "github.com/hekmekk/git-team/v2/src/shared/activation/interface"
	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/v2/src/shared/config/interface"
	state "github.com/hekmekk/git-team/v2/src/shared/state/interface"
)

// Dependencies the dependencies of the status Policy module
type Dependencies struct {
	StateReader         state.Reader
	ConfigReader        config.Reader
	ActivationValidator activation.Validator
	StateAsJson         bool
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply show the current status of git-team
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	cfg, cfgReadErr := deps.ConfigReader.Read()
	if cfgReadErr != nil {
		return StateRetrievalFailed{Reason: fmt.Errorf("failed to read config: %s", cfgReadErr)}
	}

	activationScope := cfg.ActivationScope

	if activationScope == activationscope.RepoLocal && !deps.ActivationValidator.IsInsideAGitRepository() {
		return StateRetrievalFailed{Reason: fmt.Errorf("failed to get status with activation-scope=%s: not inside a git repository", activationScope)}
	}

	retState, stateRepositoryQueryErr := deps.StateReader.Query(cfg.ActivationScope)
	if stateRepositoryQueryErr != nil {
		return StateRetrievalFailed{Reason: fmt.Errorf("failed to query current state: %s", stateRepositoryQueryErr)}
	}

	return StateRetrievalSucceeded{State: retState, StateAsJson: deps.StateAsJson}
}
