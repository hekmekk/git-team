package interfaceadapter

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

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

	effects := MapAddEventToEffects(add.AssignmentSucceeded{Alias: alias, Coauthor: coauthor})

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

	effects := MapAddEventToEffects(add.AssignmentFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsAssignmentAborted(t *testing.T) {
	expectedEffects := []effects.Effect{
		effects.NewPrintMessage("Nothing changed."),
		effects.NewExitOk(),
	}

	effects := MapAddEventToEffects(add.AssignmentAborted{Alias: "", ExistingCoauthor: "", ReplacingCoauthor: ""})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsUnknownEvent(t *testing.T) {
	expectedEffects := []effects.Effect{}

	effects := MapAddEventToEffects("UNKNOWN_EVENT")

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}
