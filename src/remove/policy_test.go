package remove

import (
	"errors"
	"reflect"
	"testing"
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

func TestRmShouldRemoveTheAssignmentWhenTryingToRemoveANonExistingAlias(t *testing.T) {
	alias := "mr"

	remove := func(alias string) error {
		return errors.New("exit status 5")
	}

	expectedEvent := DeAllocationSucceeded{Alias: alias}

	event := Policy{Dependencies{GitRemoveAlias: remove}, DeAllocationRequest{Alias: &alias}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestRmShouldFailBecauseUnderlyingGitRemoveFails(t *testing.T) {
	alias := "mr"

	err := errors.New("git remove DeAllocationRequest failed")

	remove := func(alias string) error {
		return err
	}

	expectedEvent := DeAllocationFailed{Reason: err}

	event := Policy{Dependencies{GitRemoveAlias: remove}, DeAllocationRequest{Alias: &alias}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
