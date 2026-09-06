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
	"fmt"
	"reflect"
	"testing"

	gitconfigerror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
)

func TestRmShouldRemoveTheAssignment(t *testing.T) {
	alias := "mr"

	remove := func(alias string) error {
		return nil
	}

	expectedEvent := DeAllocationSucceeded{Alias: alias}

	event := Policy{Dependencies{GitRemoveAlias: remove}, DeAllocationRequest{Alias: &alias}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestRmShouldNotRemoveTheAssignmentWhenTryingToRemoveANonExistingAlias(t *testing.T) {
	alias := "mr"

	remove := func(alias string) error {
		return gitconfigerror.ErrTryingToUnsetAnOptionWhichDoesNotExist
	}

	expectedEvent := DeAllocationFailed{Reason: fmt.Errorf("no such alias: '%s'", alias)}

	event := Policy{Dependencies{GitRemoveAlias: remove}, DeAllocationRequest{Alias: &alias}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestRmShouldFailBecauseUnderlyingGitRemoveFails(t *testing.T) {
	alias := "mr"

	remove := func(alias string) error {
		return gitconfigerror.ErrConfigFileCannotBeWritten
	}

	expectedEvent := DeAllocationFailed{Reason: fmt.Errorf("failed to remove alias: %s", gitconfigerror.ErrConfigFileCannotBeWritten)}

	event := Policy{Dependencies{GitRemoveAlias: remove}, DeAllocationRequest{Alias: &alias}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
