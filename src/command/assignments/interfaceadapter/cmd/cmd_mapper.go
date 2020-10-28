package assignmentscmdadapter

import (
	"github.com/urfave/cli/v2"

	addcmdadapter "github.com/hekmekk/git-team/src/command/assignments/add/interfaceadapter/cmd"
	listcmdadapter "github.com/hekmekk/git-team/src/command/assignments/list/interfaceadapter/cmd"
	removecmdadapter "github.com/hekmekk/git-team/src/command/assignments/remove/interfaceadapter/cmd"
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
