package removeeventadapter

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
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/v2/src/command/assignments/remove"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
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

	expectedEffect := effects.NewExitErrMsg(err)

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
