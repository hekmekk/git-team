package add

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/effects"
)

func TestAddSucceeds(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	add := func(alias, coAuthor string) error {
		return nil
	}

	msg := fmt.Sprintf("Alias '%s' -> '%s' has been added.", alias, coAuthor)

	expectedEffects := []effects.Effect{effects.NewPrintMessage(msg), effects.NewExitOk()}

	effects := Run(Dependencies{AddGitAlias: add}, Args{Alias: &alias, Coauthor: &coAuthor})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestAddFailsDueToProvidedCoauthorNotPassingSanityCheck(t *testing.T) {
	alias := "mr"
	coAuthor := "INVALID COAUTHOR"

	add := func(alias, coAuthor string) error {
		return nil
	}

	expectedEffects := []effects.Effect{effects.NewPrintErr(errors.New("Not a valid coauthor: INVALID COAUTHOR")), effects.NewExitErr()}

	effects := Run(Dependencies{AddGitAlias: add}, Args{Alias: &alias, Coauthor: &coAuthor})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestAddFailsBecauseUnderlyingGitAddFails(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	expectedErr := errors.New("git add failed")

	add := func(alias, coAuthor string) error {
		return expectedErr
	}

	expectedEffects := []effects.Effect{effects.NewPrintErr(expectedErr), effects.NewExitErr()}

	effects := Run(Dependencies{AddGitAlias: add}, Args{Alias: &alias, Coauthor: &coAuthor})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}
