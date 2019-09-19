package status

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/core/state"
)

func TestStatusShouldBeRetrieved(t *testing.T) {
	currState := state.NewStateDisabled()

	deps := Dependencies{
		StateRepositoryQuery: func() (state.State, error) { return currState, nil },
	}

	expectedEvent := StateRetrievalSucceeded{State: currState}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestStatusShouldNotBeRetrieved(t *testing.T) {
	err := errors.New("failed to retrieve state")

	deps := Dependencies{
		StateRepositoryQuery: func() (state.State, error) { return state.State{}, err },
	}

	expectedEvent := StateRetrievalFailed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
