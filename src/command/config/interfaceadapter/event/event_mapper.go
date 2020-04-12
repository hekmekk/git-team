package configeventadapter

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/fatih/color"
	configevents "github.com/hekmekk/git-team/src/command/config/events"
	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
)

// MapEventToEffects convert config events to effects for the cli
func MapEventToEffects(event events.Event) []effects.Effect {
	switch evt := event.(type) {
	case configevents.RetrievalSucceeded:
		return []effects.Effect{
			effects.NewPrintMessage(toString(evt.Config)),
			effects.NewExitOk(),
		}
	case configevents.RetrievalFailed:
		return []effects.Effect{
			effects.NewPrintErr(evt.Reason),
			effects.NewExitErr(),
		}
	default:
		return []effects.Effect{}
	}
}

func toString(cfg config.Config) string {
	roProperties := make(map[string]string)
	rwProperties := make(map[string]string)
	roProperties["commit-template-path"] = cfg.Ro.GitTeamCommitTemplatePath
	roProperties["hooks-path"] = cfg.Ro.GitTeamHooksPath
	rwProperties["activation-scope"] = cfg.Rw.ActivationScope.String()

	var properties []string
	for k, v := range roProperties {
		properties = append(properties, fmt.Sprintf("[ro] %s: %s", k, v))
	}

	for k, v := range rwProperties {
		properties = append(properties, fmt.Sprintf("[rw] %s: %s", k, v))
	}

	sort.Strings(properties)

	var buffer bytes.Buffer
	buffer.WriteString(color.New(color.FgBlue).Add(color.Bold).Sprint("config"))
	for _, property := range properties {
		buffer.WriteString(color.WhiteString("\n─ %s", property))
	}

	return buffer.String()
}
