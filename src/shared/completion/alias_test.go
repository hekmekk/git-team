package completion

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

func TestComplete(t *testing.T) {
	t.Parallel()

	cases := []struct {
		selectedAliases          []string
		expectedRemainingAliases []string
	}{
		{[]string{}, []string{"alias1", "alias2", "alias3"}},
		{[]string{"alias1"}, []string{"alias2", "alias3"}},
		{[]string{"alias1", "alias2"}, []string{"alias3"}},
		{[]string{"alias1", "alias2", "alias3"}, []string{}},
	}

	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return map[string]string{
				"team.alias.alias1": "Mr. Noujz <noujz@mr.se>",
				"team.alias.alias2": "Mrs. Noujz <noujz@mrs.se>",
				"team.alias.alias3": "Mrs. Very Noujz <very-noujz@mrs.se>",
			}, nil
		},
	}

	aliasShellCompletion := NewAliasShellCompletion(gitConfigReader)

	for _, caseLoopVar := range cases {
		selectedAliases := caseLoopVar.selectedAliases
		expectedRemainingAliases := caseLoopVar.expectedRemainingAliases

		t.Run(fmt.Sprintf("selectedAliases: %s expectedRemainingAliases: %s", selectedAliases, expectedRemainingAliases), func(t *testing.T) {
			t.Parallel()

			remainingAliases := aliasShellCompletion.Complete(selectedAliases)

			if !reflect.DeepEqual(expectedRemainingAliases, remainingAliases) {
				t.Errorf("expected: %s, actual: %s", expectedRemainingAliases, remainingAliases)
				t.Fail()
			}
		})
	}
}

func TestCompleteWhenNoAssignmentsExists(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return map[string]string{}, nil
		},
	}

	expectedRemainingAliases := []string{}

	aliasShellCompletion := NewAliasShellCompletion(gitConfigReader)

	remainingAliases := aliasShellCompletion.Complete([]string{})

	if !reflect.DeepEqual(expectedRemainingAliases, remainingAliases) {
		t.Errorf("expected: %s, actual: %s", expectedRemainingAliases, remainingAliases)
		t.Fail()
	}
}

func TestCompletewhenLookingUpAssignmentsFails(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return map[string]string{}, errors.New("any kind of error")
		},
	}

	expectedRemainingAliases := []string{}

	aliasShellCompletion := NewAliasShellCompletion(gitConfigReader)

	remainingAliases := aliasShellCompletion.Complete([]string{})

	if !reflect.DeepEqual(expectedRemainingAliases, remainingAliases) {
		t.Errorf("expected: %s, actual: %s", expectedRemainingAliases, remainingAliases)
		t.Fail()
	}
}
