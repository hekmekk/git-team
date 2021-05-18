package removeeventadapter

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/command/assignments/remove"
	"github.com/hekmekk/git-team/src/shared/cliadapter/effects"
)

func TestMapEventToEffectDeAllocationSucceeded(t *testing.T) {
	alias := "mr"
	msg := fmt.Sprintf("Assignment removed: '%s'", alias)

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(remove.DeAllocationSucceeded{Alias: alias})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectDeAllocationFailed(t *testing.T) {
	err := errors.New("failure")

	expectedEffect := effects.NewExitErr(err)

	effect := MapEventToEffect(remove.DeAllocationFailed{Reason: err})

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
