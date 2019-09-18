package listeventadapter

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/core/assignment"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/list"
)

func TestMapEventToEffectsRetrievalSucceededSorted(t *testing.T) {
	assignments := []assignment.Assignment{
		assignment.Assignment{Alias: "alias2", Coauthor: "coauthor2"},
		assignment.Assignment{Alias: "alias1", Coauthor: "coauthor1"},
		assignment.Assignment{Alias: "alias3", Coauthor: "coauthor3"},
	}

	msg := fmt.Sprintf("Assignments:\n------------\n'alias1' -> 'coauthor1'\n'alias2' -> 'coauthor2'\n'alias3' -> 'coauthor3'")

	expectedEffects := []effects.Effect{
		effects.NewPrintMessage(msg),
		effects.NewExitOk(),
	}

	effects := MapEventToEffects(list.RetrievalSucceeded{Assignments: assignments})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsRetrievalSucceededEmpty(t *testing.T) {
	assignments := []assignment.Assignment{}

	msg := fmt.Sprintf("Assignments:\n------------")

	expectedEffects := []effects.Effect{
		effects.NewPrintMessage(msg),
		effects.NewExitOk(),
	}

	effects := MapEventToEffects(list.RetrievalSucceeded{Assignments: assignments})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsRetrievalFailed(t *testing.T) {
	err := errors.New("failure")

	expectedEffects := []effects.Effect{
		effects.NewPrintErr(err),
		effects.NewExitErr(),
	}

	effects := MapEventToEffects(list.RetrievalFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsUnknownEvent(t *testing.T) {
	expectedEffects := []effects.Effect{}

	effects := MapEventToEffects("UNKNOWN_EVENT")

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}
