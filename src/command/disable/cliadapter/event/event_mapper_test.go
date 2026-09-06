package disableeventadapter

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

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
