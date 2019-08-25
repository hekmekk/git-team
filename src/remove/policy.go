package remove

import (
	"github.com/hekmekk/git-team/src/core/events"
	giterror "github.com/hekmekk/git-team/src/core/gitconfig/error"
)

// DeAllocationRequest remove an alias -> coauthor assignment
type DeAllocationRequest struct {
	Alias *string
}

// Dependencies the dependencies of the remove Policy module
type Dependencies struct {
	GitRemoveAlias func(string) error
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
	Req  DeAllocationRequest
}

// Apply remove an alias -> coauthor assignment
func (policy Policy) Apply() events.Event {
	deps := policy.Deps
	req := policy.Req

	alias := *req.Alias

	err := deps.GitRemoveAlias(alias)
	if err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return DeAllocationFailed{Reason: err}
	}

	return DeAllocationSucceeded{Alias: alias}
}
