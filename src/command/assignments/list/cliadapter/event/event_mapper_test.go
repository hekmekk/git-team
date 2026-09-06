package listeventadapter

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

	"github.com/hekmekk/git-team/v2/src/command/assignments/list"
	"github.com/hekmekk/git-team/v2/src/core/assignment"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
)

func TestMapEventToEffectRetrievalSucceededSorted(t *testing.T) {
	assignments := []assignment.Assignment{
		assignment.Assignment{Alias: "alias200", Coauthor: "coauthor2"},
		assignment.Assignment{Alias: "alias1", Coauthor: "coauthor1"},
		assignment.Assignment{Alias: "alias3", Coauthor: "coauthor3"},
	}

	msg := fmt.Sprintf("Assignments\n─ alias1   →  coauthor1\n─ alias200 →  coauthor2\n─ alias3   →  coauthor3")

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(list.RetrievalSucceeded{Assignments: assignments, OnlyAlias: false})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectRetrievalSucceededSortedWithoutCoauthors(t *testing.T) {
	assignments := []assignment.Assignment{
		assignment.Assignment{Alias: "alias200", Coauthor: "coauthor2"},
		assignment.Assignment{Alias: "alias1", Coauthor: "coauthor1"},
		assignment.Assignment{Alias: "alias3", Coauthor: "coauthor3"},
	}

	msg := fmt.Sprintf("alias1\nalias200\nalias3")

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(list.RetrievalSucceeded{Assignments: assignments, OnlyAlias: true})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectRetrievalSucceededEmpty(t *testing.T) {
	assignments := []assignment.Assignment{}

	msg := fmt.Sprintf("No assignments")

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(list.RetrievalSucceeded{Assignments: assignments, OnlyAlias: false})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectRetrievalSucceededEmptyWithoutCoauthors(t *testing.T) {
	assignments := []assignment.Assignment{}

	msg := fmt.Sprintf("")

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(list.RetrievalSucceeded{Assignments: assignments, OnlyAlias: true})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectRetrievalFailed(t *testing.T) {
	err := errors.New("failure")

	expectedEffect := effects.NewExitErrMsg(err)

	effect := MapEventToEffect(list.RetrievalFailed{Reason: err})

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
