package remove

import (
	"github.com/hekmekk/git-team/src/core/events"
	giterror "github.com/hekmekk/git-team/src/gitconfig/error"
)

// DeAllocationRequest remove an alias -> coauthor assignment
type DeAllocationRequest struct {
	Alias *string
}

// Dependencies the dependencies of the Runner
type Dependencies struct {
	GitRemoveAlias func(string) error
}

// Policy the policy to apply
type Policy struct {
	deps Dependencies
	req  DeAllocationRequest
}

// NewPolicy constructor of Policy
func NewPolicy(deps Dependencies, req DeAllocationRequest) Policy {
	return Policy{deps, req}
}

// Apply remove an alias -> coauthor assignment
func (policy Policy) Apply() events.Event {
	deps := policy.deps
	req := policy.req

	alias := *req.Alias

	err := deps.GitRemoveAlias(alias)
	if err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return DeAllocationFailed{Reason: err}
	}

	return DeAllocationSucceeded{Alias: alias}
}
