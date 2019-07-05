package handler

import (
	"errors"
	"testing"
)

func TestAddSucceeds(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	add := func(alias, coAuthor string) error {
		return nil
	}

	handleAdd := RunAddCommand(add)
	aliasAdded, err := handleAdd(alias, coAuthor)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := AliasAdded{Alias: alias, CoAuthor: coAuthor}

	if aliasAdded != expected {
		t.Errorf("was: %s, expected: %s", aliasAdded, expected)
		t.Fail()
	}
}

func TestAddFailsBecauseUnderlyingGitAddFails(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	add := func(alias, coAuthor string) error {
		return errors.New("git add failed")
	}

	handleAdd := RunAddCommand(add)
	aliasAdded, err := handleAdd(alias, coAuthor)

	if err == nil {
		t.Fail()
	}

	expected := AliasAdded{}

	if aliasAdded != expected {
		t.Errorf("was: %s, expected: %s", aliasAdded, expected)
		t.Fail()
	}
}
