package disableeventadapter

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/v2/src/command/disable"
	status "github.com/hekmekk/git-team/v2/src/command/status"
	"github.com/hekmekk/git-team/v2/src/core/events"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
	state "github.com/hekmekk/git-team/v2/src/shared/state/entity"
)

type statusPolicyMock struct {
	apply func() events.Event
}

func (mock statusPolicyMock) Apply() events.Event {
	return mock.apply()
}

func TestMapEventToEffectSucceeded(t *testing.T) {
	expectedEffect := effects.NewExitOkMsg("git-team disabled")

	statusPolicy := &statusPolicyMock{
		apply: func() events.Event {
			return status.StateRetrievalSucceeded{State: state.NewStateDisabled()}
		},
	}

	effect := MapEventToEffectFactory(statusPolicy)(disable.Succeeded{})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectSucceededButQueryStatusFailed(t *testing.T) {
	err := errors.New("query status failure")

	expectedEffect := effects.NewExitErrMsg(err)

	statusPolicy := &statusPolicyMock{
		apply: func() events.Event {
			return status.StateRetrievalFailed{Reason: err}
		},
	}

	effect := MapEventToEffectFactory(statusPolicy)(disable.Succeeded{})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectFailed(t *testing.T) {
	err := errors.New("disable failure")

	expectedEffect := effects.NewExitErrMsg(err)

	effect := MapEventToEffectFactory(nil)(disable.Failed{Reason: err})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectUnknownEvent(t *testing.T) {
	expectedEffect := effects.NewExitOk()

	effect := MapEventToEffectFactory(nil)("UNKNOWN_EVENT")

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}
