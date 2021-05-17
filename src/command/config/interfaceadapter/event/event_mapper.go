package configeventadapter

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/fatih/color"

	configevents "github.com/hekmekk/git-team/src/command/config/events"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
)

// MapEventToEffect convert config events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case configevents.RetrievalSucceeded:
		return effects.NewExitOkMsg(toString(evt.Config))
	case configevents.RetrievalFailed:
		return effects.NewExitErr(evt.Reason)
	case configevents.SettingModificationSucceeded:
		return effects.NewExitOkMsg(color.CyanString(fmt.Sprintf("Configuration updated: '%s' → '%s'", evt.Key, evt.Value)))
	case configevents.SettingModificationFailed:
		return effects.NewExitErr(evt.Reason)
	case configevents.ReadingSingleSettingNotYetImplemented:
		return effects.NewExitErr(errors.New("Reading a single setting has not yet been implemented"))
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
		buffer.WriteString(color.WhiteString("\n─ %s", property))
	}

	return buffer.String()
}
