package statuseventadapter

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/command/status"
	"github.com/hekmekk/git-team/src/shared/cli/effects"
	state "github.com/hekmekk/git-team/src/shared/state/entity"
)

func TestMapEventToEffectStateRetrievalSucceededEnabled(t *testing.T) {
	msg := "git-team enabled\n\nco-authors\n─ Mr. Noujz <noujz@mr.se>\n─ Mrs. Noujz <noujz@mrs.se>"
	state := state.NewStateEnabled([]string{"Mrs. Noujz <noujz@mrs.se>", "Mr. Noujz <noujz@mr.se>"}, "/previous/hooks/path")

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(status.StateRetrievalSucceeded{State: state, StateAsJson: false})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectStateRetrievalSucceededEnabledJson(t *testing.T) {
	msg := `{"status":"enabled","coAuthors":["Mrs. Noujz <noujz@mrs.se>","Mr. Noujz <noujz@mr.se>"],"previousHooksPath":"/previous/hooks/path"}` + "\n"
	state := state.NewStateEnabled([]string{"Mrs. Noujz <noujz@mrs.se>", "Mr. Noujz <noujz@mr.se>"}, "/previous/hooks/path")

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(status.StateRetrievalSucceeded{State: state, StateAsJson: true})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectStateRetrievalSucceededDisabled(t *testing.T) {
	msg := "git-team disabled"
	state := state.NewStateDisabled()

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(status.StateRetrievalSucceeded{State: state, StateAsJson: false})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectStateRetrievalSucceededDisabledJson(t *testing.T) {
	msg := `{"status":"disabled","coAuthors":[],"previousHooksPath":""}` + "\n"
	state := state.NewStateDisabled()

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(status.StateRetrievalSucceeded{State: state, StateAsJson: true})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectStateRetrievalFailed(t *testing.T) {
	err := errors.New("failure")

	expectedEffect := effects.NewExitErrMsg(err)

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
