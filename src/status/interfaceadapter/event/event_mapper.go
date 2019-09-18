package statuseventadapter

import (
	"bytes"
	"sort"

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/state"
	"github.com/hekmekk/git-team/src/status"
)

// MapEventToEffects convert status events to effects for the cli
func MapEventToEffects(event events.Event) []effects.Effect {
	switch evt := event.(type) {
	case status.StateRetrievalSucceeded:
		return []effects.Effect{
			effects.NewPrintMessage(toString(evt.State)),
			effects.NewExitOk(),
		}
	case status.StateRetrievalFailed:
		return []effects.Effect{
			effects.NewPrintErr(evt.Reason),
			effects.NewExitErr(),
		}
	default:
		return []effects.Effect{}
	}
}

const msgTemplate string = "git-team %s."

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
				buffer.WriteString(color.WhiteString("\n- %s", coauthor))
			}
		}
	}

	return buffer.String()
}
