package addeventadapter

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/command/assignments/add"
	"github.com/hekmekk/git-team/src/core/effects"
)

func TestMapEventToEffectsAssignmentSucceeded(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"
	msg := fmt.Sprintf("Assignment added: '%s' →  '%s'", alias, coauthor)

	expectedEffects := []effects.Effect{
		effects.NewPrintMessage(msg),
		effects.NewExitOk(),
	}

	effects := MapEventToEffects(add.AssignmentSucceeded{Alias: alias, Coauthor: coauthor})

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

	effects := MapEventToEffects(add.AssignmentFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsAssignmentAborted(t *testing.T) {
	expectedEffects := []effects.Effect{
		effects.NewPrintMessage("Nothing changed"),
		effects.NewExitOk(),
	}

	effects := MapEventToEffects(add.AssignmentAborted{Alias: "", ExistingCoauthor: "", ReplacingCoauthor: ""})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsUnknownEvent(t *testing.T) {
	expectedEffects := []effects.Effect{}

	effects := MapEventToEffects("UNKNOWN_EVENT")

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}
