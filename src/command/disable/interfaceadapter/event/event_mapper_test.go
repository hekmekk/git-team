package disableeventadapter

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/command/disable"
	status "github.com/hekmekk/git-team/src/command/status"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/state"
)

type statusPolicyMock struct {
	apply func() events.Event
}

func (mock statusPolicyMock) Apply() events.Event {
	return mock.apply()
}

func TestMapEventToEffectsSucceeded(t *testing.T) {
	expectedEffects := []effects.Effect{
		effects.NewPrintMessage("git-team disabled"),
		effects.NewExitOk(),
	}

	statusPolicy := &statusPolicyMock{
		apply: func() events.Event {
			return status.StateRetrievalSucceeded{State: state.NewStateDisabled()}
		},
	}

	effects := MapEventToEffectsFactory(statusPolicy)(disable.Succeeded{})

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

	statusPolicy := &statusPolicyMock{
		apply: func() events.Event {
			return status.StateRetrievalFailed{Reason: err}
		},
	}

	effects := MapEventToEffectsFactory(statusPolicy)(disable.Succeeded{})

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
