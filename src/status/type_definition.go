package status

import (
	"bytes"

	"github.com/fatih/color"
)

type teamStatus string

const (
	enabled     teamStatus = "enabled"
	disabled    teamStatus = "disabled"
	msgTemplate string     = "git-team %s."
)

type state struct {
	Status    teamStatus `toml:"status"`
	Coauthors []string   `toml:"co-authors"`
}

func (state state) ToString() string {
	var buffer bytes.Buffer
	switch state.Status {
	case enabled:
		buffer.WriteString(color.GreenString(msgTemplate, enabled))
		coauthors := state.Coauthors
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
