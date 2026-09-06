package commandadapter

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
	"github.com/hekmekk/git-team/v2/src/core/events"
	"github.com/hekmekk/git-team/v2/src/core/policy"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
)

// Run apply a policy, convert the resulting event to an effect and run that effect
func Run(policy policy.Policy, eventMapper func(events.Event) effects.Effect) error {
	return ApplyPolicy(policy, eventMapper).Run()
}

// ApplyPolicy apply a policy and convert the resulting event to an effect
func ApplyPolicy(policy policy.Policy, eventMapper func(events.Event) effects.Effect) effects.Effect {
	return eventMapper(policy.Apply())
}
