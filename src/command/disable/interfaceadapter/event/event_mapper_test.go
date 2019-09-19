package disableeventadapter

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/command/disable"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/state"
)

func TestMapEventToEffectsSucceeded(t *testing.T) {
	expectedEffects := []effects.Effect{
		effects.NewPrintMessage("git-team disabled"),
		effects.NewExitOk(),
	}

	queryStatus := func() (state.State, error) {
		return state.NewStateDisabled(), nil
	}

	effects := MapEventToEffectsFactory(queryStatus)(disable.Succeeded{})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsSucceededButQueryStatusFailed(t *testing.T) {
	err := errors.New("query status failure")

	expectedEffects := []effects.Effect{
		effects.NewPrintErr(err),
		effects.NewExitErr(),
	}

	queryStatus := func() (state.State, error) {
		return state.State{}, err
	}

	effects := MapEventToEffectsFactory(queryStatus)(disable.Succeeded{})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsFailed(t *testing.T) {
	err := errors.New("disable failure")

	expectedEffects := []effects.Effect{
		effects.NewPrintErr(err),
		effects.NewExitErr(),
	}

	effects := MapEventToEffectsFactory(nil)(disable.Failed{Reason: err})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsUnknownEvent(t *testing.T) {
	expectedEffects := []effects.Effect{}

	effects := MapEventToEffectsFactory(nil)("UNKNOWN_EVENT")

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}
