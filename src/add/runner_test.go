package add

import (
	"errors"
	"reflect"
	"testing"
)

func TestAddSucceeds(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"

	add := func(alias, coauthor string) error {
		return nil
	}

	expectedEvent := AssignmentSucceeded{Alias: alias, Coauthor: coauthor}

	event := Run(Dependencies{AddGitAlias: add}, Args{Alias: &alias, Coauthor: &coauthor})

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddFailsDueToProvidedCoauthorNotPassingSanityCheck(t *testing.T) {
	alias := "mr"
	coauthor := "INVALID COAUTHOR"

	err := errors.New("Not a valid coauthor: INVALID COAUTHOR")

	add := func(alias, coAuthor string) error {
		return nil
	}

	expectedEvent := AssignmentFailed{Reason: err}

	event := Run(Dependencies{AddGitAlias: add}, Args{Alias: &alias, Coauthor: &coauthor})

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddFailsBecauseUnderlyingGitAddFails(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"

	err := errors.New("git add failed")

	add := func(alias, coauthor string) error {
		return err
	}

	expectedEvent := AssignmentFailed{Reason: err}

	event := Run(Dependencies{AddGitAlias: add}, Args{Alias: &alias, Coauthor: &coauthor})

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
