package listeventadapter

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
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/v2/src/command/assignments/list"
	"github.com/hekmekk/git-team/v2/src/core/assignment"
	"github.com/hekmekk/git-team/v2/src/core/events"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
)

// MapEventToEffect convert list events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case list.RetrievalSucceeded:
		return effects.NewExitOkMsg(toString(evt))
	case list.RetrievalFailed:
		return effects.NewExitErrMsg(evt.Reason)
	default:
		return effects.NewExitOk()
	}
}

func toString(evt list.RetrievalSucceeded) string {
	sorted := evt.Assignments
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Alias < sorted[j].Alias })

	if evt.OnlyAlias {
		aliases := []string{}
		for _, assignment := range sorted {
			aliases = append(aliases, assignment.Alias)
		}
		return strings.Join(aliases, "\n")
	}

	return toStringWithCoauthors(sorted)
}

func toStringWithCoauthors(assignments []assignment.Assignment) string {
	maxAliasLength := 0
	for _, assignment := range assignments {
		currAliasLength := len(assignment.Alias)
		if currAliasLength > maxAliasLength {
			maxAliasLength = currAliasLength
		}
	}

	var buffer bytes.Buffer

	if len(assignments) == 0 {
		buffer.WriteString(color.New(color.FgBlue).Add(color.Bold).Sprint("No assignments"))
		return buffer.String()
	}

	buffer.WriteString(color.New(color.FgBlue).Add(color.Bold).Sprint("Assignments"))
	for _, assignment := range assignments {
		buffer.WriteString(fmt.Sprintf("\n─ %-[1]*s →  %s", maxAliasLength, assignment.Alias, assignment.Coauthor))
	}

	return buffer.String()
}
