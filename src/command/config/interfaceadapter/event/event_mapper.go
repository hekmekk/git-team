package configeventadapter

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/fatih/color"

	config "github.com/hekmekk/git-team/src/command/config/entity/config"
	configevents "github.com/hekmekk/git-team/src/command/config/events"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
)

// MapEventToEffects convert config events to effects for the cli
func MapEventToEffects(event events.Event) []effects.Effect {
	switch evt := event.(type) {
	case configevents.RetrievalSucceeded:
		return []effects.Effect{
			effects.NewPrintMessage(toString(evt.Config)),
			effects.NewExitOk(),
		}
	case configevents.RetrievalFailed:
		return []effects.Effect{
			effects.NewPrintErr(evt.Reason),
			effects.NewExitErr(),
		}
	case configevents.SettingModificationSucceeded:
		return []effects.Effect{
			effects.NewPrintMessage(color.CyanString(fmt.Sprintf("Configuration updated: '%s' → '%s'", evt.Key, evt.Value))),
			effects.NewExitOk(),
		}
	case configevents.SettingModificationFailed:
		return []effects.Effect{
			effects.NewPrintErr(evt.Reason),
			effects.NewExitErr(),
		}
	case configevents.ReadingSingleSettingNotYetImplemented:
		return []effects.Effect{
			effects.NewPrintErr(errors.New("Reading a single setting has not yet been implemented")),
			effects.NewExitErr(),
		}
	default:
		return []effects.Effect{}
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
