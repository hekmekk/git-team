package disableeventadapter

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
	"github.com/hekmekk/git-team/v2/src/command/disable"
	statuseventadapter "github.com/hekmekk/git-team/v2/src/command/status/cliadapter/event"
	"github.com/hekmekk/git-team/v2/src/core/events"
	"github.com/hekmekk/git-team/v2/src/core/policy"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
)

// MapEventToEffectFactory convert disable events to effects for the cli
func MapEventToEffectFactory(statusPolicy policy.Policy) func(events.Event) effects.Effect {
	return func(event events.Event) effects.Effect {
		switch evt := event.(type) {
		case disable.Succeeded:
			return statuseventadapter.MapEventToEffect(statusPolicy.Apply())
		case disable.Failed:
			return effects.NewExitErrMsg(evt.Reason)
		default:
			return effects.NewExitOk()
		}
	}
}
