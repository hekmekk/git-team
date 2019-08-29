package list

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/core/assignment"
)

func TestListShouldReturnTheAvailableAssignments(t *testing.T) {
	aliasCoauthorMap := map[string]string{
		"alias1": "coauthor1",
		"alias2": "coauthor2",
	}

	deps := Dependencies{
		GitGetAssignments: func() (map[string]string, error) { return aliasCoauthorMap, nil },
	}

	assignmentsA := []assignment.Assignment{
		assignment.Assignment{Alias: "alias1", Coauthor: "coauthor1"},
		assignment.Assignment{Alias: "alias2", Coauthor: "coauthor2"},
	}
	expectedEventA := RetrievalSucceeded{Assignments: assignmentsA}

	assignmentsB := []assignment.Assignment{
		assignment.Assignment{Alias: "alias2", Coauthor: "coauthor2"},
		assignment.Assignment{Alias: "alias1", Coauthor: "coauthor1"},
	}
	expectedEventB := RetrievalSucceeded{Assignments: assignmentsB}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEventA, event) && !reflect.DeepEqual(expectedEventB, event) {
		t.Errorf("expected: %s, got: %s", expectedEventA, event)
		t.Fail()
	}
}

func TestListShouldReturnFailure(t *testing.T) {
	err := errors.New("failed to get assignments")

	deps := Dependencies{
		GitGetAssignments: func() (map[string]string, error) { return map[string]string{}, err },
	}

	expectedEvent := RetrievalFailed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
