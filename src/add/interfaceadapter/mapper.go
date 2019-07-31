package interfaceadapter

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/src/add"
	"github.com/hekmekk/git-team/src/effects"
)

// MapAddEventToEffects convert assignment events to effects for the cli
func MapAddEventToEffects(event interface{}) []effects.Effect {
	switch x := event.(type) {
	case add.AssignmentSucceeded:
		alias := x.Alias
		coauthor := x.Coauthor
		return []effects.Effect{
			effects.NewPrintMessage(color.GreenString(fmt.Sprintf("Alias '%s' -> '%s' has been added.", alias, coauthor))),
			effects.NewExitOk(),
		}

	case add.AssignmentFailed:
		return []effects.Effect{
			effects.NewPrintErr(x.Reason),
			effects.NewExitErr(),
		}
	default:
		return []effects.Effect{}
	}
}
