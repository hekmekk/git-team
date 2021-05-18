package listeventadapter

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/command/assignments/list"
	"github.com/hekmekk/git-team/src/core/assignment"
	"github.com/hekmekk/git-team/src/shared/cliadapter/effects"
)

func TestMapEventToEffectRetrievalSucceededSorted(t *testing.T) {
	assignments := []assignment.Assignment{
		assignment.Assignment{Alias: "alias200", Coauthor: "coauthor2"},
		assignment.Assignment{Alias: "alias1", Coauthor: "coauthor1"},
		assignment.Assignment{Alias: "alias3", Coauthor: "coauthor3"},
	}

	msg := fmt.Sprintf("Assignments\n─ alias1   →  coauthor1\n─ alias200 →  coauthor2\n─ alias3   →  coauthor3")

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(list.RetrievalSucceeded{Assignments: assignments})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectRetrievalSucceededEmpty(t *testing.T) {
	assignments := []assignment.Assignment{}

	msg := fmt.Sprintf("No assignments")

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(list.RetrievalSucceeded{Assignments: assignments})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectRetrievalFailed(t *testing.T) {
	err := errors.New("failure")

	expectedEffect := effects.NewExitErr(err)

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
