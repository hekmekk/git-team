package statuseventadapter

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

	"github.com/hekmekk/git-team/v2/src/command/status"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
	state "github.com/hekmekk/git-team/v2/src/shared/state/entity"
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
	msg := `{"status":"enabled","coAuthors":["Mrs. Noujz <noujz@mrs.se>","Mr. Noujz <noujz@mr.se>"]}`
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
	msg := `{"status":"disabled","coAuthors":[]}`
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
