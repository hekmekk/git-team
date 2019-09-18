package listeventadapter

import (
	"bytes"
	"sort"

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/src/core/assignment"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/list"
)

// MapEventToEffects convert list events to effects for the cli
func MapEventToEffects(event events.Event) []effects.Effect {
	switch evt := event.(type) {
	case list.RetrievalSucceeded:
		return []effects.Effect{
			effects.NewPrintMessage(toString(evt.Assignments)),
			effects.NewExitOk(),
		}
	case list.RetrievalFailed:
		return []effects.Effect{
			effects.NewPrintErr(evt.Reason),
			effects.NewExitErr(),
		}
	default:
		return []effects.Effect{}
	}
}

func toString(assignments []assignment.Assignment) string {
	sorted := assignments
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Alias < sorted[j].Alias })

	maxAliasLength := 0
	for _, assignment := range sorted {
		currAliasLength := len(assignment.Alias)
		if currAliasLength > maxAliasLength {
			maxAliasLength = currAliasLength
		}
	}

	var buffer bytes.Buffer

	buffer.WriteString(color.New(color.FgBlue).Add(color.Bold).Sprint("assignments"))

	for _, assgnmnt := range sorted {
		buffer.WriteString(color.WhiteString("\n- %-[1]*s -> '%s'", maxAliasLength, assgnmnt.Alias, assgnmnt.Coauthor))
	}

	return buffer.String()
}
