package completioncmdadapter

import (
	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/src/command/completion/bash"
	commandadapter "github.com/hekmekk/git-team/src/shared/cli/commandadapter"
	"github.com/hekmekk/git-team/src/shared/cli/effects"
)

func Command() *cli.Command {
	bashCommand := &cli.Command{
		Name:        "bash",
		Usage:       "Bash completion",
		Description: "Source with bash to get auto completion",
		Action: func(c *cli.Context) error {
			return commandadapter.RunEffect(effects.NewExitOkMsg(bash.Script))
		},
	}

	return &cli.Command{
		Name:        "completion",
		Usage:       "Shell completion",
		Description: "Source the output of this command to get auto completion",
		Subcommands: []*cli.Command{
			bashCommand,
		},
	}
}