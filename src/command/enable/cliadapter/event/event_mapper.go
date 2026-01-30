package enableeventadapter

import (
	"bytes"
	"errors"
	"strings"

	"github.com/hekmekk/git-team/v2/src/command/enable"
	statuseventadapter "github.com/hekmekk/git-team/v2/src/command/status/cliadapter/event"
	"github.com/hekmekk/git-team/v2/src/core/events"
	"github.com/hekmekk/git-team/v2/src/core/policy"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
)

// MapEventToEffectFactory convert enable events to effects for the cli
func MapEventToEffectFactory(statusPolicy policy.Policy) func(events.Event) effects.Effect {
	return func(event events.Event) effects.Effect {
		switch evt := event.(type) {
		case enable.Succeeded, enable.Aborted:
			return statuseventadapter.MapEventToEffect(statusPolicy.Apply())
		case enable.Failed:
			return effects.NewExitErrMsg(foldErrors(evt.Reason))
		default:
			return effects.NewExitOk()
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
