package disableeventadapter

import (
	"github.com/hekmekk/git-team/src/command/disable"
	statuseventadapter "github.com/hekmekk/git-team/src/command/status/cliadapter/event"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/policy"
	"github.com/hekmekk/git-team/src/shared/cli/effects"
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
