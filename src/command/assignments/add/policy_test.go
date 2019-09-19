package add

import (
	"errors"
	"reflect"
	"testing"
)

func TestAddShouldMakeTheNewAssignment(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return nil },
		GitResolveAlias:     func(alias string) (string, error) { return "", errors.New("No such alias") },
		GitAddAlias:         func(alias, coauthor string) error { return nil },
		GetAnswerFromUser:   func(string) (string, error) { return "", nil },
	}

	expectedEvent := AssignmentSucceeded{Alias: alias, Coauthor: coauthor}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &coauthor}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddShouldFailDueToProvidedCoauthorNotPassingSanityCheck(t *testing.T) {
	alias := "mr"
	coauthor := "INVALID COAUTHOR"

	err := errors.New("Not a valid coauthor: INVALID COAUTHOR")

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return err },
	}

	expectedEvent := AssignmentFailed{Reason: err}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &coauthor}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddShouldNotOverrideTheOriginalAssignment(t *testing.T) {
	alias := "mr"
	existingCoauthor := "Mr. Green <green@mr.se>"
	replacingCoauthor := "Mr. Noujz <noujz@mr.se>"

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return nil },
		GitResolveAlias:     func(alias string) (string, error) { return existingCoauthor, nil },
		GitAddAlias:         func(alias, coauthor string) error { return nil },
		GetAnswerFromUser:   func(string) (string, error) { return "", nil },
	}

	expectedEvent := AssignmentAborted{Alias: alias, ExistingCoauthor: existingCoauthor, ReplacingCoauthor: replacingCoauthor}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &replacingCoauthor}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddShouldOverrideTheOriginalAssignment(t *testing.T) {
	alias := "mr"
	existingCoauthor := "Mr. Green <green@mr.se>"
	replacingCoauthor := "Mr. Noujz <noujz@mr.se>"

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return nil },
		GitResolveAlias:     func(alias string) (string, error) { return existingCoauthor, nil },
		GitAddAlias:         func(alias, coauthor string) error { return nil },
		GetAnswerFromUser:   func(string) (string, error) { return "y", nil },
	}

	expectedEvent := AssignmentSucceeded{Alias: alias, Coauthor: replacingCoauthor}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &replacingCoauthor}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddShouldFailBecauseUnderlyingGitAddFails(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"

	err := errors.New("git add failed")

	deps := Dependencies{
		SanityCheckCoauthor: func(coauthor string) error { return nil },
		GitResolveAlias:     func(alias string) (string, error) { return "", errors.New("No such alias") },
		GitAddAlias:         func(alias, coauthor string) error { return err },
	}

	expectedEvent := AssignmentFailed{Reason: err}

	event := Policy{deps, AssignmentRequest{Alias: &alias, Coauthor: &coauthor}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
