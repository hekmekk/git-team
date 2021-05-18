package commandadapter

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/policy"
)

// RunUrFave apply a policy and convert the resulting event(s) to effects
func RunUrFave(policy policy.Policy, adapter func(events.Event) effects.Effect) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		effect := adapter(policy.Apply())
		return RunEffect(effect)
	}
}

// RunEffect derive the action to be performed based on the more abstract effect
func RunEffect(effect effects.Effect) error {
	switch e := effect.(type) {
	case effects.ExitOk:
		fmt.Println(e.Message())
		return nil
	case effects.ExitWarn:
		fmt.Println(color.YellowString(e.Message()))
		return nil
	case effects.ExitErr:
		return cli.NewExitError(color.RedString(e.Message()), 1)
	default:
		return cli.NewExitError(color.RedString("undefined behavior encountered"), 42)
	}
}
