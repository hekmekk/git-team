package configeventadapter

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
	"errors"
	"fmt"
	"sort"

	"github.com/fatih/color"

	configevents "github.com/hekmekk/git-team/v2/src/command/config/events"
	"github.com/hekmekk/git-team/v2/src/core/events"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
	config "github.com/hekmekk/git-team/v2/src/shared/config/entity/config"
)

// MapEventToEffect convert config events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case configevents.RetrievalSucceeded:
		return effects.NewExitOkMsg(toString(evt.Config))
	case configevents.RetrievalFailed:
		return effects.NewExitErrMsg(evt.Reason)
	case configevents.SettingModificationSucceeded:
		return effects.NewExitOkMsg(color.CyanString(fmt.Sprintf("Configuration updated: '%s' → '%s'", evt.Key, evt.Value)))
	case configevents.SettingModificationFailed:
		return effects.NewExitErrMsg(evt.Reason)
	case configevents.ReadingSingleSettingNotYetImplemented:
		return effects.NewExitErrMsg(errors.New("Reading a single setting has not yet been implemented"))
	default:
		return effects.NewExitOk()
	}
}

func toString(cfg config.Config) string {
	properties := make(map[string]string)
	properties["activation-scope"] = cfg.ActivationScope.String()

	var propertyStrings []string

	for k, v := range properties {
		propertyStrings = append(propertyStrings, fmt.Sprintf("%s: %s", k, v))
	}

	sort.Strings(propertyStrings)

	var buffer bytes.Buffer
	buffer.WriteString(color.New(color.FgBlue).Add(color.Bold).Sprint("config"))
	for _, property := range propertyStrings {
		buffer.WriteString(fmt.Sprintf("\n─ %s", property))
	}

	return buffer.String()
}
