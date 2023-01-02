package statuseventadapter

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/src/command/status"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/shared/cli/effects"
	state "github.com/hekmekk/git-team/src/shared/state/entity"
)

// MapEventToEffect convert status events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case status.StateRetrievalSucceeded:
		return effects.NewExitOkMsg(toString(evt.State))
	case status.StateRetrievalFailed:
		return effects.NewExitErrMsg(evt.Reason)
	default:
		return effects.NewExitOk()
	}
}

const msgTemplate string = "git-team %s"

func toString(theState state.State) string {
	var buffer bytes.Buffer
	buffer.WriteString(color.CyanString(msgTemplate, theState.Status))
	if theState.IsEnabled() {
		coauthors := theState.Coauthors
		sort.Strings(coauthors)
		if len(coauthors) > 0 {
			buffer.WriteString("\n\n")
			buffer.WriteString(color.New(color.FgBlue).Add(color.Bold).Sprint("co-authors"))
			for _, coauthor := range coauthors {
				buffer.WriteString(fmt.Sprintf("\n─ %s", coauthor))
			}
		}
	}

	return buffer.String()
}
