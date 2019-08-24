package disableeventadapter

import (
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/disable"
	"github.com/hekmekk/git-team/src/status"
)

// MapEventToEffectsFactory convert assignment events to effects for the cli
func MapEventToEffectsFactory(queryStatus func() (status.State, error)) func(events.Event) []effects.Effect {
	return func(event events.Event) []effects.Effect {
		switch evt := event.(type) {
		case disable.Succeeded:
			// TODO: replace this with status policy application + event mapper
			status, err := queryStatus()
			if err != nil {
				return []effects.Effect{
					effects.NewPrintErr(err),
					effects.NewExitErr(),
				}
			}

			return []effects.Effect{
				effects.NewPrintMessage(status.ToString()),
				effects.NewExitOk(),
			}
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
