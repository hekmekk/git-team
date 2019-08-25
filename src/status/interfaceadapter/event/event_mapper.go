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
	if theState.IsEnabled() {
		buffer.WriteString(color.GreenString(msgTemplate, theState.Status))
		coauthors := theState.Coauthors
		sort.Strings(coauthors)
		if len(coauthors) > 0 {
			buffer.WriteString("\n\n")
			blackBold := color.New(color.FgBlack).Add(color.Bold)
			buffer.WriteString(blackBold.Sprintln("Co-authors:"))
			buffer.WriteString(blackBold.Sprint("-----------"))
			for _, coauthor := range coauthors {
				buffer.WriteString(color.MagentaString("\n%s", coauthor))
			}
		}
	} else {
		buffer.WriteString(color.RedString(msgTemplate, theState.Status))
	}

	return buffer.String()
}
