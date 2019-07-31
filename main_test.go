package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"testing/quick"

	"github.com/hekmekk/git-team/src/add"
	"github.com/hekmekk/git-team/src/effects"
)

func TestMapEventToEffectsAssignmentSucceeded(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"
	msg := fmt.Sprintf("Alias '%s' -> '%s' has been added.", alias, coauthor)

	expectedEffects := []effects.Effect{
		effects.NewPrintMessage(msg),
		effects.NewExitOk(),
	}

	effects := mapEventToEffects(add.AssignmentSucceeded{Alias: alias, Coauthor: coauthor})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsAssignmentFailed(t *testing.T) {
	err := errors.New("failure")

	expectedEffects := []effects.Effect{
		effects.NewPrintErr(err),
		effects.NewExitErr(),
	}

	effects := mapEventToEffects(add.AssignmentFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestFoldErrors(t *testing.T) {
	errPrefix := errors.New("_prefix_")
	errSuffix := errors.New("_suffix_")

	// Note: It is more than twice as slow with this predicate approach... Maybe revert to direct inline calls
	isNotNil := func(err error) bool { return err != nil }
	hasProperPrefix := func(err error) bool { return strings.HasPrefix(err.Error(), errPrefix.Error()) }
	hasProperSuffix := func(err error) bool { return strings.HasSuffix(err.Error(), errSuffix.Error()) }

	errorsGen := func(msg string) bool {
		generatedErr := errors.New(msg)
		errs := []error{errPrefix, generatedErr, errSuffix}

		if foldedErr := foldErrors(errs); isNotNil(foldedErr) && hasProperPrefix(foldedErr) && hasProperSuffix(foldedErr) {
			return true
		}
		return false
	}

	if err := quick.Check(errorsGen, nil); err != nil {
		t.Error(err)
	}
}
