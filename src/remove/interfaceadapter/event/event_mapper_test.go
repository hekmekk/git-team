package removeeventadapter

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/effects"
	"github.com/hekmekk/git-team/src/remove"
)

func TestMapEventToEffectsDeAllocationSucceeded(t *testing.T) {
	alias := "mr"
	msg := fmt.Sprintf("Alias '%s' has been removed.", alias)

	expectedEffects := []effects.Effect{
		effects.NewPrintMessage(msg),
		effects.NewExitOk(),
	}

	effects := MapRemoveEventToEffects(remove.DeAllocationSucceeded{Alias: alias})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsDeAllocationFailed(t *testing.T) {
	err := errors.New("failure")

	expectedEffects := []effects.Effect{
		effects.NewPrintErr(err),
		effects.NewExitErr(),
	}

	effects := MapRemoveEventToEffects(remove.DeAllocationFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsUnknownEvent(t *testing.T) {
	expectedEffects := []effects.Effect{}

	effects := MapRemoveEventToEffects("UNKNOWN_EVENT")

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}
