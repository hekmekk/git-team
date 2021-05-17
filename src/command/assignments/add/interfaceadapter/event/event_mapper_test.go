package addeventadapter

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/command/assignments/add"
	"github.com/hekmekk/git-team/src/core/effects"
)

func TestMapEventToEffectAssignmentSucceeded(t *testing.T) {
	alias := "mr"
	coauthor := "Mr. Noujz <noujz@mr.se>"
	msg := fmt.Sprintf("Assignment added: '%s' →  '%s'", alias, coauthor)

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(add.AssignmentSucceeded{Alias: alias, Coauthor: coauthor})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectAssignmentFailed(t *testing.T) {
	err := errors.New("failure")

	expectedEffect := effects.NewExitErr(err)

	effect := MapEventToEffect(add.AssignmentFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectAssignmentAborted(t *testing.T) {
	expectedEffect := effects.NewExitOkMsg("Nothing changed")
	effect := MapEventToEffect(add.AssignmentAborted{Alias: "", ExistingCoauthor: "", ReplacingCoauthor: ""})

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
