package removeeventadapter

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/hekmekk/git-team/src/effects"
	"github.com/hekmekk/git-team/src/remove"
)

// MapEventToEffects convert deallocation events to effects for the cli
func MapEventToEffects(event interface{}) []effects.Effect {
	switch evt := event.(type) {
	case remove.DeAllocationSucceeded:
		return []effects.Effect{
			effects.NewPrintMessage(color.RedString(fmt.Sprintf("Alias '%s' has been removed.", evt.Alias))),
			effects.NewExitOk(),
		}
	case remove.DeAllocationFailed:
		return []effects.Effect{
			effects.NewPrintErr(evt.Reason),
			effects.NewExitErr(),
		}
	default:
		return []effects.Effect{}
	}
}
