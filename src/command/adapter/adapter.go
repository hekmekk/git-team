package commandadapter

import (
	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/policy"
)

// RunUrFave apply a policy and convert the resulting event(s) to effects
func RunUrFave(policy policy.Policy, adapter func(events.Event) []effects.Effect) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		for _, effect := range adapter(policy.Apply()) {
			effect.Run()
		}
		return nil
	}
}
