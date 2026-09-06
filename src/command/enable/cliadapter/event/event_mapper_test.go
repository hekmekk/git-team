package enableeventadapter

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
	"strings"
	"testing"
	"testing/quick"

	"github.com/hekmekk/git-team/v2/src/command/enable"
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

func TestFoldErrors(t *testing.T) {
	errPrefix := errors.New("_prefix_")
	errSuffix := errors.New("_suffix_")

	// Note: It is more than twice as slow with this predicate approach... Maybe revert to direct inline calls
	isNotNil := func(err error) bool { return err != nil }
	hasProperPrefix := func(err error) bool { return strings.HasPrefix(err.Error(), errPrefix.Error()) }
	hasProperSuffix := func(err error) bool { return strings.HasSuffix(err.Error(), errSuffix.Error()) }

	errorsGen := func(msg string) bool {
		generatedErr := errors.New(msg)
		errs := []error{errPrefix, generatedErr, errSuffix}

		if foldedErr := foldErrors(errs); isNotNil(foldedErr) && hasProperPrefix(foldedErr) && hasProperSuffix(foldedErr) {
			return true
		}
		return false
	}

	if err := quick.Check(errorsGen, nil); err != nil {
		t.Error(err)
	}
}

func TestMapEventToEffectSucceeded(t *testing.T) {
	coauthors := []string{"A Coauthor"}
	expectedEffect := effects.NewExitOkMsg(fmt.Sprintf("git-team enabled\n\nco-authors\n─ A Coauthor"))

	statusPolicy := &statusPolicyMock{
		apply: func() events.Event {
			return status.StateRetrievalSucceeded{State: state.NewStateEnabled(coauthors, "/previous/hooks/path")}
		},
	}

	effect := MapEventToEffectFactory(statusPolicy)(enable.Succeeded{})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectAborted(t *testing.T) {
	expectedEffect := effects.NewExitOkMsg("git-team disabled")

	statusPolicy := &statusPolicyMock{
		apply: func() events.Event {
			return status.StateRetrievalSucceeded{State: state.NewStateDisabled()}
		},
	}

	effect := MapEventToEffectFactory(statusPolicy)(enable.Aborted{})

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

	effect := MapEventToEffectFactory(statusPolicy)(enable.Succeeded{})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectFailed(t *testing.T) {
	err := errors.New("enable failure")

	expectedEffect := effects.NewExitErrMsg(err)

	effect := MapEventToEffectFactory(nil)(enable.Failed{Reason: []error{err}})

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
