package commandadapter

import (
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/policy"
	"github.com/hekmekk/git-team/src/shared/cli/effects"
)

// Run apply a policy, convert the resulting event to an effect and run that effect
func Run(policy policy.Policy, eventMapper func(events.Event) effects.Effect) error {
	return ApplyPolicy(policy, eventMapper).Run()
}

// ApplyPolicy apply a policy and convert the resulting event to an effect
func ApplyPolicy(policy policy.Policy, eventMapper func(events.Event) effects.Effect) effects.Effect {
	return eventMapper(policy.Apply())
}
