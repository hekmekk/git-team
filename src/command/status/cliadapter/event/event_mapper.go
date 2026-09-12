package statuseventadapter

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
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/hekmekk/git-team/v2/src/command/status"
	"github.com/hekmekk/git-team/v2/src/core/events"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
	state "github.com/hekmekk/git-team/v2/src/shared/state/entity"
)

// MapEventToEffect convert status events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case status.StateRetrievalSucceeded:
		if evt.StateAsJson {
			return toJson(evt.State)
		}
		return effects.NewExitOkMsg(toString(evt.State))
	case status.StateRetrievalFailed:
		return effects.NewExitErrMsg(evt.Reason)
	default:
		return effects.NewExitOk()
	}
}

func toString(theState state.State) string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("git-team %s", theState.Status))
	if theState.IsEnabled() {
		coauthors := theState.Coauthors
		sort.Strings(coauthors)
		if len(coauthors) > 0 {
			buffer.WriteString("\n\n")
			buffer.WriteString("co-authors")
			for _, coauthor := range coauthors {
				buffer.WriteString(fmt.Sprintf("\n─ %s", coauthor))
			}
		}
	}

	return buffer.String()
}

func toJson(theState state.State) effects.Effect {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)

	jsonStructure := struct {
		Status    string   `json:"status"`
		Coauthors []string `json:"coAuthors"`
	}{
		Status:    string(theState.Status),
		Coauthors: theState.Coauthors,
	}

	err := encoder.Encode(jsonStructure)
	if err != nil {
		return effects.NewExitErrMsg(err)
	}

	jsonStringWithoutTrailingNewline := strings.TrimSuffix(buffer.String(), "\n")
	return effects.NewExitOkMsg(jsonStringWithoutTrailingNewline)
}
