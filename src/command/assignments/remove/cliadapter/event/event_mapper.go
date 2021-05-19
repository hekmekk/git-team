package removeeventadapter

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/src/command/assignments/remove"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/shared/cli/effects"
)

// MapEventToEffect convert deallocation events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case remove.DeAllocationSucceeded:
		return effects.NewExitOkMsg(color.CyanString(fmt.Sprintf("Assignment removed: '%s'", evt.Alias)))
	case remove.DeAllocationFailed:
		return effects.NewExitErrMsg(evt.Reason)
	default:
		return effects.NewExitOk()
	}
}
