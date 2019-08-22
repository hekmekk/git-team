package remove

import (
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

// Apply remove an alias -> coauthor assignment
func Apply(deps Dependencies, req DeAllocationRequest) interface{} {
	alias := *req.Alias

	err := deps.GitRemoveAlias(alias)
	if err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return DeAllocationFailed{Reason: err}
	}

	return DeAllocationSucceeded{Alias: alias}
}
