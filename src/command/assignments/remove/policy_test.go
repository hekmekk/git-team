package remove

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
