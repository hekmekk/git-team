package enableeventadapter

import (
	"bytes"
	"errors"
	"strings"

	"github.com/hekmekk/git-team/src/command/enable"
	"github.com/hekmekk/git-team/src/command/status"
	"github.com/hekmekk/git-team/src/command/status/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/state"
)

// MapEventToEffectsFactory convert enable events to effects for the cli
func MapEventToEffectsFactory(queryStatus func() (state.State, error)) func(events.Event) []effects.Effect {
	return func(event events.Event) []effects.Effect {
		switch evt := event.(type) {
		case enable.Succeeded, enable.Aborted:
			return statuseventadapter.MapEventToEffects(status.Policy{Deps: status.Dependencies{StateRepositoryQuery: queryStatus}}.Apply())
		case enable.Failed:
			return []effects.Effect{
				effects.NewPrintErr(foldErrors(evt.Reason)),
				effects.NewExitErr(),
			}
		default:
			return []effects.Effect{}
		}
	}
}

func foldErrors(validationErrors []error) error {
	var buffer bytes.Buffer
	for _, err := range validationErrors {
		buffer.WriteString(err.Error())
		buffer.WriteString("; ")
	}
	return errors.New(strings.TrimRight(buffer.String(), "; "))
}
