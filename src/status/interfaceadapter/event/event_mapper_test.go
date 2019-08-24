package statuseventadapter

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/state"
	"github.com/hekmekk/git-team/src/status"
)

func TestMapEventToEffectsStateRetrievalSucceededEnabled(t *testing.T) {
	msg := "git-team enabled.\n\nCo-authors:\n-----------\nMr. Noujz <noujz@mr.se>\nMrs. Noujz <noujz@mrs.se>"
	state := state.NewStateEnabled([]string{"Mrs. Noujz <noujz@mrs.se>", "Mr. Noujz <noujz@mr.se>"})

	expectedEffects := []effects.Effect{
		effects.NewPrintMessage(msg),
		effects.NewExitOk(),
	}

	effects := MapEventToEffects(status.StateRetrievalSucceeded{State: state})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsStateRetrievalSucceededDisabled(t *testing.T) {
	msg := "git-team disabled."
	state := state.NewStateDisabled()

	expectedEffects := []effects.Effect{
		effects.NewPrintMessage(msg),
		effects.NewExitOk(),
	}

	effects := MapEventToEffects(status.StateRetrievalSucceeded{State: state})

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

	effects := MapEventToEffects(status.StateRetrievalFailed{Reason: err})

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
