package assignmentscmdadapter

import (
	"github.com/urfave/cli/v2"

	addcmdadapter "github.com/hekmekk/git-team/v2/src/command/assignments/add/cliadapter/cmd"
	listcmdadapter "github.com/hekmekk/git-team/v2/src/command/assignments/list/cliadapter/cmd"
	removecmdadapter "github.com/hekmekk/git-team/v2/src/command/assignments/remove/cliadapter/cmd"
)

// Command the assignments command
func Command() *cli.Command {
	return &cli.Command{
		Name:   "assignments",
		Usage:  "Manage your alias -> co-author assignments",
		Action: listcmdadapter.Command().Action,
		Subcommands: []*cli.Command{
			addcmdadapter.Command(),
			listcmdadapter.Command(),
			removecmdadapter.Command(),
		},
	}
}
