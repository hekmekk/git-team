package addeventadapter

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/src/command/assignments/add"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/shared/cli/effects"
)

// MapEventToEffect convert assignment events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case add.AssignmentSucceeded:
		return effects.NewExitOkMsg(color.CyanString(fmt.Sprintf("Assignment added: '%s' →  '%s'", evt.Alias, evt.Coauthor)))
	case add.AssignmentFailed:
		return effects.NewExitErr(evt.Reason)
	case add.AssignmentAborted:
		return effects.NewExitOkMsg("Nothing changed")
	default:
		return effects.NewExitOk()
	}
}
