package disableeventadapter

import (
	"github.com/hekmekk/git-team/src/command/disable"
	statuseventadapter "github.com/hekmekk/git-team/src/command/status/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/policy"
)

// MapEventToEffectsFactory convert disable events to effects for the cli
func MapEventToEffectsFactory(statusPolicy policy.Policy) func(events.Event) []effects.Effect {
	return func(event events.Event) []effects.Effect {
		switch evt := event.(type) {
		case disable.Succeeded:
			return statuseventadapter.MapEventToEffects(statusPolicy.Apply())
		case disable.Failed:
			return []effects.Effect{
				effects.NewPrintErr(evt.Reason),
				effects.NewExitErr(),
			}
		default:
			return []effects.Effect{}
		}
	}
}
