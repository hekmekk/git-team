package list

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/hekmekk/git-team/v2/src/core/assignment"
	gitconfigerror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
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

	for _, caseLoopVar := range []bool{true, false} {
		onlyAlias := caseLoopVar
		t.Run(strconv.FormatBool(onlyAlias), func(t *testing.T) {
			t.Parallel()

			req := ListRequest{
				OnlyAlias: &onlyAlias,
			}

			deps := Dependencies{
				GitConfigReader: gitConfigReader,
			}

			assignmentsA := []assignment.Assignment{
				assignment.Assignment{Alias: "alias1", Coauthor: "coauthor1"},
				assignment.Assignment{Alias: "alias2", Coauthor: "coauthor2"},
			}
			expectedEventA := RetrievalSucceeded{Assignments: assignmentsA, OnlyAlias: *req.OnlyAlias}

			assignmentsB := []assignment.Assignment{
				assignment.Assignment{Alias: "alias2", Coauthor: "coauthor2"},
				assignment.Assignment{Alias: "alias1", Coauthor: "coauthor1"},
			}
			expectedEventB := RetrievalSucceeded{Assignments: assignmentsB, OnlyAlias: *req.OnlyAlias}

			event := Policy{deps, req}.Apply()

			if !reflect.DeepEqual(expectedEventA, event) && !reflect.DeepEqual(expectedEventB, event) {
				t.Errorf("expected: %v, got: %s", expectedEventA, event)
				t.Fail()
			}
		})
	}
}

func TestListShouldReturnAnEmptyListIfGitConfigSectionIsEmpty(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return map[string]string{}, gitconfigerror.ErrSectionOrKeyIsInvalid
		},
	}

	for _, caseLoopVar := range []bool{true, false} {
		onlyAlias := caseLoopVar
		t.Run(strconv.FormatBool(onlyAlias), func(t *testing.T) {
			t.Parallel()

			req := ListRequest{
				OnlyAlias: &onlyAlias,
			}

			deps := Dependencies{
				GitConfigReader: gitConfigReader,
			}

			expectedEvent := RetrievalSucceeded{Assignments: []assignment.Assignment{}, OnlyAlias: *req.OnlyAlias}

			event := Policy{deps, req}.Apply()

			if !reflect.DeepEqual(expectedEvent, event) {
				t.Errorf("expected: %v, got: %s", expectedEvent, event)
				t.Fail()
			}
		})
	}
}

func TestListShouldReturnFailure(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return map[string]string{}, gitconfigerror.ErrTryingToUseAnInvalidRegexp
		},
	}

	onlyAlias := false

	req := ListRequest{
		OnlyAlias: &onlyAlias,
	}

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
	}

	expectedEvent := RetrievalFailed{Reason: fmt.Errorf("failed to retrieve assignments: %s", gitconfigerror.ErrTryingToUseAnInvalidRegexp)}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
