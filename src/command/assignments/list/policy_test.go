package list

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/core/assignment"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

type gitConfigReaderMock struct {
	getRegexp func(gitconfigscope.Scope, string) (map[string]string, error)
}

func (mock gitConfigReaderMock) Get(scope gitconfigscope.Scope, key string) (string, error) {
	return "", nil
}

func (mock gitConfigReaderMock) GetAll(scope gitconfigscope.Scope, key string) ([]string, error) {
	return []string{}, nil
}

func (mock gitConfigReaderMock) GetRegexp(scope gitconfigscope.Scope, pattern string) (map[string]string, error) {
	return mock.getRegexp(scope, pattern)
}

func (mock gitConfigReaderMock) List(scope gitconfigscope.Scope) (map[string]string, error) {
	return nil, nil
}

func TestListShouldReturnTheAvailableAssignments(t *testing.T) {
	aliasCoauthorMap := map[string]string{
		"team.alias.alias1": "coauthor1",
		"team.alias.alias2": "coauthor2",
	}

	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return aliasCoauthorMap, nil
		},
	}

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
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

func TestListShouldReturnAnEmptyListIfGitConfigSectionIsEmpty(t *testing.T) {
	err := errors.New("exit status 1")

	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return map[string]string{}, err
		},
	}

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
	}

	expectedEvent := RetrievalSucceeded{Assignments: []assignment.Assignment{}}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestListShouldReturnFailure(t *testing.T) {
	err := errors.New("failed to get assignments")

	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return map[string]string{}, err
		},
	}

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
	}

	expectedEvent := RetrievalFailed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
