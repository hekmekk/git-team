package remove

import (
	"errors"
	"fmt"

	"github.com/hekmekk/git-team/src/core/events"
	gitconfigerror "github.com/hekmekk/git-team/src/shared/gitconfig/error"
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
	if err != nil {
		if errors.Is(err, gitconfigerror.ErrTryingToUnsetAnOptionWhichDoesNotExist) {
			return DeAllocationFailed{Reason: fmt.Errorf("no such alias: '%s'", alias)}
		}

		return DeAllocationFailed{Reason: fmt.Errorf("failed to remove alias: %s", err)}
	}

	return DeAllocationSucceeded{Alias: alias}
}
