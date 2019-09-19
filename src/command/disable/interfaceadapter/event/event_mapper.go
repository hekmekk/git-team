package disableeventadapter

import (
	"github.com/hekmekk/git-team/src/command/disable"
	"github.com/hekmekk/git-team/src/command/status"
	"github.com/hekmekk/git-team/src/command/status/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/state"
)

// MapEventToEffectsFactory convert disable events to effects for the cli
func MapEventToEffectsFactory(queryStatus func() (state.State, error)) func(events.Event) []effects.Effect {
	return func(event events.Event) []effects.Effect {
		switch evt := event.(type) {
		case disable.Succeeded:
			return statuseventadapter.MapEventToEffects(status.Policy{Deps: status.Dependencies{StateRepositoryQuery: queryStatus}}.Apply())
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
