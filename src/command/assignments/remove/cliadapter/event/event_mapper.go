package removeeventadapter

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

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/v2/src/command/assignments/remove"
	"github.com/hekmekk/git-team/v2/src/core/events"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
)

// MapEventToEffect convert deallocation events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case remove.DeAllocationSucceeded:
		return effects.NewExitOkMsg(color.CyanString(fmt.Sprintf("Assignment removed: '%s'", evt.Alias)))
	case remove.DeAllocationFailed:
		return effects.NewExitErrMsg(evt.Reason)
	default:
		return effects.NewExitOk()
	}
}
