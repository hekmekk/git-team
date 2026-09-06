package add

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
	"fmt"
	"strings"

	"github.com/hekmekk/git-team/v2/src/core/events"
)

// AssignmentRequest which coauthor to assign to the alias
type AssignmentRequest struct {
	Alias         *string
	Coauthor      *string
	ForceOverride *bool
	KeepExisting  *bool
}

// Dependencies the dependencies of the add Policy module
type Dependencies struct {
	SanityCheckCoauthor func(string) error
	GitAddAlias         func(string, string) error
	GitResolveAlias     func(string) (string, error)
	GetAnswerFromUser   func(string) (string, error)
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
	Req  AssignmentRequest
}

const (
	y   string = "y"
	yes string = "yes"
)

type assignmentReplacementStrategy int

const (
	Override assignmentReplacementStrategy = iota
	KeepExisting
	Ask
)

// Apply assign a co-author to an alias
func (policy Policy) Apply() events.Event {
	req := policy.Req
	deps := policy.Deps

	alias := *req.Alias
	coauthor := *req.Coauthor

	if checkErr := deps.SanityCheckCoauthor(coauthor); checkErr != nil {
		return AssignmentFailed{Reason: checkErr}
	}

	forceOverride := *req.ForceOverride
	keepExisting := *req.KeepExisting

	assignmentReplacementStrategy := deriveAssignmentReplacementStrategy(forceOverride, keepExisting)

	isAssignmentExisting := false

	existingCoauthor, resolveErr := deps.GitResolveAlias(alias)
	if resolveErr == nil {
		isAssignmentExisting = true
	} else {
		isAssignmentExisting = false
	}

	shouldAddAssignment := true

	if isAssignmentExisting && assignmentReplacementStrategy == Ask {
		choice, err := shouldAssignmentBeOverridden(deps, alias, existingCoauthor, coauthor)
		if err != nil {
			return AssignmentFailed{Reason: fmt.Errorf("failed to retrieve answer from user: %s", err)}
		}
		shouldAddAssignment = choice
	} else {
		if assignmentReplacementStrategy == Override {
			shouldAddAssignment = true
		}

		if assignmentReplacementStrategy == KeepExisting {
			shouldAddAssignment = false
		}
	}

	if !shouldAddAssignment {
		return AssignmentAborted{}
	}

	err := deps.GitAddAlias(alias, coauthor)
	if err != nil {
		return AssignmentFailed{Reason: fmt.Errorf("failed to add alias: %s", err)}
	}

	return AssignmentSucceeded{Alias: alias, Coauthor: coauthor}
}

func deriveAssignmentReplacementStrategy(forceOverride bool, keepExisting bool) assignmentReplacementStrategy {
	if forceOverride == keepExisting {
		return Ask
	}

	if forceOverride {
		return Override
	}

	return KeepExisting
}

func shouldAssignmentBeOverridden(deps Dependencies, alias, existingCoauthor, replacingCoauthor string) (bool, error) {
	question := fmt.Sprintf("Assignment '%s' →  '%s' exists already. Override with '%s'? [N/y] ", alias, existingCoauthor, replacingCoauthor)

	answer, err := deps.GetAnswerFromUser(question)
	if err != nil {
		return false, err
	}

	answer = strings.ToLower(strings.TrimSpace(strings.TrimRight(answer, "\n")))
	switch answer {
	case y, yes:
		return true, nil
	default:
		return false, nil
	}
}
