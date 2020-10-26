package commandadapter

import (
	"github.com/urfave/cli/v2"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/policy"
)

// CommandRoot kingpin application or command
type CommandRoot interface {
	Command(string, string) *kingpin.CmdClause
}

// RunUrFave apply a policy and convert the resulting event(s) to effects
func RunUrFave(policy policy.Policy, adapter func(events.Event) []effects.Effect) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		for _, effect := range adapter(policy.Apply()) {
			effect.Run()
		}
		return nil
	}
}

// Run apply a policy and convert the resulting event(s) to effects
func Run(policy policy.Policy, adapter func(events.Event) []effects.Effect) func(c *kingpin.ParseContext) error {
	return func(c *kingpin.ParseContext) error {
		for _, effect := range adapter(policy.Apply()) {
			effect.Run()
		}
		return nil
	}
}
