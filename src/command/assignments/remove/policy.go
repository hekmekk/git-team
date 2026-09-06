package remove

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"errors"
	"fmt"

	"github.com/hekmekk/git-team/v2/src/core/events"
	gitconfigerror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
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
