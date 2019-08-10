package add

import (
	"errors"
	"reflect"
	"testing"
)

func TestAddSucceeds(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"

	deps := Dependencies{
		GitResolveAlias: func(alias string) (string, error) { return "", errors.New("no such alias") },
		AddGitAlias:     func(alias, coauthor string) error { return nil },
	}

	expectedEvent := AssignmentSucceeded{Alias: alias, Coauthor: coauthor}

	event := Apply(deps, Args{Alias: &alias, Coauthor: &coauthor})

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddShouldCheckForAnExistingAliasFirst(t *testing.T) {
	alias := "mr"
	existingCoauthor := "Mr. Green <green@mr.se>"
	replacingCoauthor := "Mr. Noujz <noujz@mr.se>"

	deps := Dependencies{
		GitResolveAlias: func(alias string) (string, error) { return existingCoauthor, nil },
		AddGitAlias:     func(alias, coauthor string) error { return nil },
	}

	expectedEvent := AssignmentAborted{Alias: alias, ExistingCoauthor: existingCoauthor, ReplacingCoauthor: replacingCoauthor}

	event := Apply(deps, Args{Alias: &alias, Coauthor: &replacingCoauthor})

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddFailsDueToProvidedCoauthorNotPassingSanityCheck(t *testing.T) {
	alias := "mr"
	coauthor := "INVALID COAUTHOR"

	err := errors.New("Not a valid coauthor: INVALID COAUTHOR")

	deps := Dependencies{
		GitResolveAlias: func(alias string) (string, error) { return "", errors.New("no such alias") },
		AddGitAlias:     func(alias, coauthor string) error { return nil },
	}

	expectedEvent := AssignmentFailed{Reason: err}

	event := Apply(deps, Args{Alias: &alias, Coauthor: &coauthor})

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestAddFailsBecauseUnderlyingGitAddFails(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"

	err := errors.New("git add failed")

	deps := Dependencies{
		GitResolveAlias: func(alias string) (string, error) { return "", errors.New("no such alias") },
		AddGitAlias:     func(alias, coauthor string) error { return err },
	}

	expectedEvent := AssignmentFailed{Reason: err}

	event := Apply(deps, Args{Alias: &alias, Coauthor: &coauthor})

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
