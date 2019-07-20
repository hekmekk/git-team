package status

import (
	"bytes"

	"github.com/fatih/color"
)

type teamStatus string

const (
	enabled  teamStatus = "enabled"
	disabled teamStatus = "disabled"
)

type state struct {
	Status    teamStatus `toml:"status"`
	Coauthors []string   `toml:"co-authors"`
}

func toString(state state) string {
	var buffer bytes.Buffer
	switch state.Status {
	case enabled:
		buffer.WriteString(color.GreenString("git-team %s.", enabled))
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
		buffer.WriteString(color.RedString("git-team %s", disabled))
	}

	return buffer.String()
}
