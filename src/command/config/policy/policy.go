package policy

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

	configevents "github.com/hekmekk/git-team/v2/src/command/config/events"
	"github.com/hekmekk/git-team/v2/src/core/events"
	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/v2/src/shared/config/interface"
)

// Request defines which config setting to modify or if the config should just be displayed
type Request struct {
	Key   *string
	Value *string
}

// Dependencies the dependencies of the config Policy module
type Dependencies struct {
	ConfigWriter config.Writer
	ConfigReader config.Reader
}

// Policy the policy to apply
type Policy struct {
	Req  Request
	Deps Dependencies
}

// Apply Edit or show configuration
func (policy Policy) Apply() events.Event {
	deps := policy.Deps
	req := policy.Req
	keyPtr := req.Key
	valuePtr := req.Value

	if (keyPtr == nil || *keyPtr == "") && (valuePtr == nil || *valuePtr == "") {
		cfg, err := deps.ConfigReader.Read()
		if err != nil {
			return configevents.RetrievalFailed{Reason: fmt.Errorf("failed to read config: %s", err)}
		}

		return configevents.RetrievalSucceeded{Config: cfg}
	}

	if keyPtr == nil || *keyPtr == "" || valuePtr == nil || *valuePtr == "" {
		return configevents.ReadingSingleSettingNotYetImplemented{}
	}

	key := *keyPtr
	value := *valuePtr

	if key != "activation-scope" {
		return configevents.SettingModificationFailed{Reason: fmt.Errorf("unknown setting '%s'", key)}
	}

	desiredScope := activationscope.FromString(value)
	if desiredScope == activationscope.Unknown {
		return configevents.SettingModificationFailed{Reason: fmt.Errorf("unknown activation-scope '%s'", value)}
	}

	if err := deps.ConfigWriter.SetActivationScope(desiredScope); err != nil {
		return configevents.SettingModificationFailed{Reason: fmt.Errorf("failed to modify setting 'activation-scope': %s", err)}
	}

	return configevents.SettingModificationSucceeded{Key: key, Value: value}
}
