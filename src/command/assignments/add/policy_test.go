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
	"errors"
	"fmt"
	"reflect"
	"testing"

	gitconfigerror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
)

func TestAddShouldMakeTheNewAssignment(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"
	forceOverride := false
	keepExisting := false

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return nil },
		GitResolveAlias:     func(alias string) (string, error) { return "", errors.New("No such alias") },
		GitAddAlias:         func(alias, coauthor string) error { return nil },
		GetAnswerFromUser:   func(string) (string, error) { return "", nil },
	}

	expectedEvent := AssignmentSucceeded{Alias: alias, Coauthor: coauthor}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &coauthor, ForceOverride: &forceOverride, KeepExisting: &keepExisting}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddShouldFailDueToProvidedCoauthorNotPassingSanityCheck(t *testing.T) {
	alias := "mr"
	coauthor := "INVALID COAUTHOR"
	forceOverride := false
	keepExisting := false

	err := errors.New("not a valid coauthor: INVALID COAUTHOR")

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return err },
	}

	expectedEvent := AssignmentFailed{Reason: err}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &coauthor, ForceOverride: &forceOverride, KeepExisting: &keepExisting}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddShouldNotOverrideTheOriginalAssignment(t *testing.T) {
	alias := "mr"
	existingCoauthor := "Mr. Green <green@mr.se>"
	replacingCoauthor := "Mr. Noujz <noujz@mr.se>"
	forceOverride := false
	keepExisting := false

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return nil },
		GitResolveAlias:     func(alias string) (string, error) { return existingCoauthor, nil },
		GitAddAlias:         func(alias, coauthor string) error { return nil },
		GetAnswerFromUser:   func(string) (string, error) { return "", nil },
	}

	expectedEvent := AssignmentAborted{}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &replacingCoauthor, ForceOverride: &forceOverride, KeepExisting: &keepExisting}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddShouldOverrideTheOriginalAssignmentAfterAskingTheUser(t *testing.T) {
	alias := "mr"
	existingCoauthor := "Mr. Green <green@mr.se>"
	replacingCoauthor := "Mr. Noujz <noujz@mr.se>"

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return nil },
		GitResolveAlias:     func(alias string) (string, error) { return existingCoauthor, nil },
		GitAddAlias:         func(alias, coauthor string) error { return nil },
		GetAnswerFromUser:   func(string) (string, error) { return "y", nil },
	}

	cases := []struct {
		forceOverride bool
		keepExisting  bool
		description   string
	}{
		{false, false, "false"},
		{true, true, "true"},
	}

	for _, caseLoopVar := range cases {
		forceOverride := caseLoopVar.forceOverride
		keepExisting := caseLoopVar.keepExisting
		description := caseLoopVar.description

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			expectedEvent := AssignmentSucceeded{Alias: alias, Coauthor: replacingCoauthor}

			event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &replacingCoauthor, ForceOverride: &forceOverride, KeepExisting: &keepExisting}}.Apply()

			if !reflect.DeepEqual(expectedEvent, event) {
				t.Errorf("expected: %s, got: %s", expectedEvent, event)
				t.Fail()
			}
		})
	}
}

func TestAddShouldOverrideTheOriginalAssignmentWhenBeingForcedTo(t *testing.T) {
	alias := "mr"
	existingCoauthor := "Mr. Green <green@mr.se>"
	replacingCoauthor := "Mr. Noujz <noujz@mr.se>"
	forceOverride := true
	keepExisting := false

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return nil },
		GitResolveAlias:     func(alias string) (string, error) { return existingCoauthor, nil },
		GitAddAlias:         func(alias, coauthor string) error { return nil },
	}

	expectedEvent := AssignmentSucceeded{Alias: alias, Coauthor: replacingCoauthor}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &replacingCoauthor, ForceOverride: &forceOverride, KeepExisting: &keepExisting}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddShouldKeepTheOriginalAssignment(t *testing.T) {
	alias := "mr"
	existingCoauthor := "Mr. Green <green@mr.se>"
	replacingCoauthor := "Mr. Noujz <noujz@mr.se>"
	forceOverride := false
	keepExisting := true

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return nil },
		GitResolveAlias:     func(alias string) (string, error) { return existingCoauthor, nil },
		GitAddAlias:         func(alias, coauthor string) error { return nil },
		GetAnswerFromUser:   func(string) (string, error) { return "y", nil }, // TODO: this should not be required
	}

	expectedEvent := AssignmentAborted{}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &replacingCoauthor, ForceOverride: &forceOverride, KeepExisting: &keepExisting}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddShouldFailBecauseUnderlyingGitAddFails(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"
	forceOverride := false
	keepExisting := false

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return nil },
		GitResolveAlias:     func(alias string) (string, error) { return "", errors.New("No such alias") },
		GitAddAlias:         func(alias, coauthor string) error { return gitconfigerror.ErrConfigFileCannotBeWritten },
	}

	expectedEvent := AssignmentFailed{Reason: fmt.Errorf("failed to add alias: %s", gitconfigerror.ErrConfigFileCannotBeWritten)}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &coauthor, ForceOverride: &forceOverride, KeepExisting: &keepExisting}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
