package listeventadapter

import (
	"bytes"
	"sort"

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/src/command/assignments/list"
	"github.com/hekmekk/git-team/src/core/assignment"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
)

// MapEventToEffect convert list events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case list.RetrievalSucceeded:
		return effects.NewExitOkMsg(toString(evt.Assignments))
	case list.RetrievalFailed:
		return effects.NewExitErr(evt.Reason)
	default:
		return effects.NewExitOk()
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

	if len(sorted) == 0 {
		buffer.WriteString(color.New(color.FgBlue).Add(color.Bold).Sprint("No assignments"))
		return buffer.String()
	}

	buffer.WriteString(color.New(color.FgBlue).Add(color.Bold).Sprint("assignments"))
	for _, assignment := range sorted {
		buffer.WriteString(color.WhiteString("\n─ %-[1]*s →  %s", maxAliasLength, assignment.Alias, assignment.Coauthor))
	}

	return buffer.String()
}
