package commandadapter

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/policy"
	"github.com/hekmekk/git-team/src/shared/cli/effects"
)

// Run apply a policy, convert the resulting event to an effect and run that effect
func Run(policy policy.Policy, eventMapper func(events.Event) effects.Effect) error {
	return RunEffect(ApplyPolicy(policy, eventMapper))
}

// ApplyPolicy apply a policy and convert the resulting event to an effect
func ApplyPolicy(policy policy.Policy, eventMapper func(events.Event) effects.Effect) effects.Effect {
	return eventMapper(policy.Apply())
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
