package handler

import (
	"errors"
	"testing"
)

func TestAddSucceeds(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	resolve := func(alias string) (string, error) {
		return "", errors.New("No alias found")
	}

	add := func(alias, coAuthor string) error {
		return nil
	}

	handleAdd := RunAddCommand(resolve, add)
	aliasAdded, err, exitCode := handleAdd(alias, coAuthor)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if exitCode != 0 {
		t.Error("Exit Code != 0")
		t.Fail()
	}

	expected := AliasAdded{Alias: alias, CoAuthor: coAuthor}

	if aliasAdded != expected {
		t.Errorf("was: %s, expected: %s", aliasAdded, expected)
		t.Fail()
	}
}

func TestAddFailsBecauseAliasIsDefined(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	resolve := func(alias string) (string, error) {
		return coAuthor, nil
	}

	add := func(alias, coAuthor string) error {
		return nil
	}

	handleAdd := RunAddCommand(resolve, add)
	aliasAdded, err, exitCode := handleAdd(alias, coAuthor)

	if err == nil {
		t.Fail()
	}

	if exitCode != 0 {
		t.Error("Exit Code != 0")
		t.Fail()
	}

	expected := AliasAdded{}

	if aliasAdded != expected {
		t.Errorf("was: %s, expected: %s", aliasAdded, expected)
		t.Fail()
	}
}

func TestAddFailsBecauseUnderlyingGitAddFails(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	resolve := func(alias string) (string, error) {
		return "", errors.New("No alias found")
	}

	add := func(alias, coAuthor string) error {
		return errors.New("git add failed")
	}

	handleAdd := RunAddCommand(resolve, add)
	aliasAdded, err, exitCode := handleAdd(alias, coAuthor)

	if err == nil {
		t.Fail()
	}

	if exitCode != -1 {
		t.Error("Exit Code != -1")
		t.Fail()
	}

	expected := AliasAdded{}

	if aliasAdded != expected {
		t.Errorf("was: %s, expected: %s", aliasAdded, expected)
		t.Fail()
	}
}
