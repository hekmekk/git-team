package status

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

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
