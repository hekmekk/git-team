package statuseventadapter

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/command/status"
	"github.com/hekmekk/git-team/src/core/effects"
	state "github.com/hekmekk/git-team/src/shared/state/entity"
)

func TestMapEventToEffectStateRetrievalSucceededEnabled(t *testing.T) {
	msg := "git-team enabled\n\nco-authors\n─ Mr. Noujz <noujz@mr.se>\n─ Mrs. Noujz <noujz@mrs.se>"
	state := state.NewStateEnabled([]string{"Mrs. Noujz <noujz@mrs.se>", "Mr. Noujz <noujz@mr.se>"})

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(status.StateRetrievalSucceeded{State: state})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectStateRetrievalSucceededDisabled(t *testing.T) {
	msg := "git-team disabled"
	state := state.NewStateDisabled()

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(status.StateRetrievalSucceeded{State: state})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectStateRetrievalFailed(t *testing.T) {
	err := errors.New("failure")

	expectedEffect := effects.NewExitErr(err)

	effect := MapEventToEffect(status.StateRetrievalFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectUnknownEvent(t *testing.T) {
	expectedEffect := effects.NewExitOk()

	effect := MapEventToEffect("UNKNOWN_EVENT")

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}
