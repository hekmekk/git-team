package addeventadapter

import (
	"fmt"

	"github.com/hekmekk/git-team/v2/src/command/assignments/add"
	"github.com/hekmekk/git-team/v2/src/core/events"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
)

// MapEventToEffect convert assignment events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case add.AssignmentSucceeded:
		return effects.NewExitOkMsg(fmt.Sprintf("Assignment added: '%s' →  '%s'", evt.Alias, evt.Coauthor))
	case add.AssignmentFailed:
		return effects.NewExitErrMsg(evt.Reason)
	case add.AssignmentAborted:
		return effects.NewExitOk()
	default:
		return effects.NewExitOk()
	}
}
