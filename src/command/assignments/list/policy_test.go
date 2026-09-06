package list

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
