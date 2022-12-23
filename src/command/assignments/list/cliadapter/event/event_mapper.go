package listeventadapter

import (
	"bytes"
	"sort"
	"strings"

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/src/command/assignments/list"
	"github.com/hekmekk/git-team/src/core/assignment"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/shared/cli/effects"
)

// MapEventToEffect convert list events to effects for the cli
func MapEventToEffect(event events.Event) effects.Effect {
	switch evt := event.(type) {
	case list.RetrievalSucceeded:
		return effects.NewExitOkMsg(toString(evt))
	case list.RetrievalFailed:
		return effects.NewExitErrMsg(evt.Reason)
	default:
		return effects.NewExitOk()
	}
}

func toString(evt list.RetrievalSucceeded) string {
	sorted := evt.Assignments
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Alias < sorted[j].Alias })

	if evt.OnlyAlias {
		aliases := []string{}
		for _, assignment := range sorted {
			aliases = append(aliases, assignment.Alias)
		}
		return strings.Join(aliases, "\n")
	}

	return toStringWithCoauthors(sorted)
}

func toStringWithCoauthors(assignments []assignment.Assignment) string {
	maxAliasLength := 0
	for _, assignment := range assignments {
		currAliasLength := len(assignment.Alias)
		if currAliasLength > maxAliasLength {
			maxAliasLength = currAliasLength
		}
	}

	var buffer bytes.Buffer

	if len(assignments) == 0 {
		buffer.WriteString(color.New(color.FgBlue).Add(color.Bold).Sprint("No assignments"))
		return buffer.String()
	}

	buffer.WriteString(color.New(color.FgBlue).Add(color.Bold).Sprint("Assignments"))
	for _, assignment := range assignments {
		buffer.WriteString(color.WhiteString("\n─ %-[1]*s →  %s", maxAliasLength, assignment.Alias, assignment.Coauthor))
	}

	return buffer.String()
}
