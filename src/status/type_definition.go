package status

import (
	"bytes"
	"sort"

	"github.com/fatih/color"
)

type teamStatus string

const (
	enabled     teamStatus = "enabled"
	disabled    teamStatus = "disabled"
	msgTemplate string     = "git-team %s."
)

// State the state of git-team
type State struct {
	Status    teamStatus `toml:"status"`
	Coauthors []string   `toml:"co-authors"`
}

// IsEnabled returns true if git-team is enabled
func (state State) IsEnabled() bool {
	return state.Status == enabled
}

// ToString the cli string representation of the current state
// TODO: move to mapper...
func (state State) ToString() string {
	var buffer bytes.Buffer
	switch state.Status {
	case enabled:
		buffer.WriteString(color.GreenString(msgTemplate, enabled))
		coauthors := state.Coauthors
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
	default:
		buffer.WriteString(color.RedString(msgTemplate, disabled))
	}

	return buffer.String()
}
