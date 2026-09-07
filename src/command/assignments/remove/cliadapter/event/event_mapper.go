package removeeventadapter

import (
	"fmt"

	"github.com/hekmekk/git-team/v2/src/command/assignments/remove"
	"github.com/hekmekk/git-team/v2/src/core/events"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
)

// MapEventToEffect convert deallocation events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case remove.DeAllocationSucceeded:
		return effects.NewExitOkMsg(fmt.Sprintf("Assignment removed: '%s'", evt.Alias))
	case remove.DeAllocationFailed:
		return effects.NewExitErrMsg(evt.Reason)
	default:
		return effects.NewExitOk()
	}
}
