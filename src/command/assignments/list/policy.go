package list

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
	"strings"

	"github.com/hekmekk/git-team/v2/src/core/assignment"
	"github.com/hekmekk/git-team/v2/src/core/events"
	giterror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// ListRequest how to show the available assignments
type ListRequest struct {
	OnlyAlias *bool
}

// Dependencies the dependencies of the list Policy module
type Dependencies struct {
	GitConfigReader gitconfig.Reader
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
	Req  ListRequest
}

// Apply show the available assignments
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	aliasCoauthorMap, err := deps.GitConfigReader.GetRegexp(gitconfigscope.Global, "team.alias")
	if err != nil && !errors.Is(err, giterror.ErrSectionOrKeyIsInvalid) {
		return RetrievalFailed{Reason: fmt.Errorf("failed to retrieve assignments: %s", err)}
	}

	assignments := []assignment.Assignment{}

	for rawAlias, coauthor := range aliasCoauthorMap {
		alias := strings.TrimPrefix(rawAlias, "team.alias.")
		assignments = append(assignments, assignment.Assignment{Alias: alias, Coauthor: coauthor})
	}

	return RetrievalSucceeded{Assignments: assignments, OnlyAlias: *policy.Req.OnlyAlias}
}
