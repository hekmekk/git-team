package gitconfig

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

	scope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

func TestListShouldReturnTheKeyValueMap(t *testing.T) {
	t.Parallel()

	key1 := "section1.key1"
	key2 := "section1.key2"
	val1 := "space separated value 1"
	val2 := "/path/to/smth"
	key3 := "section2.key1"
	val3 := "+refs/heads/*:refs/remotes/origin/*"

	expectedMap := make(map[string]string)
	expectedMap[key1] = val1
	expectedMap[key2] = val2
	expectedMap[key3] = val3

	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("%s=%s", key1, val1))
	lines = append(lines, fmt.Sprintf("%s=%s", key2, val2))
	lines = append(lines, fmt.Sprintf("%s=%s", key3, val3))

	for _, caseLoopVar := range []scope.Scope{scope.Global, scope.Local} {
		gitConfigScope := caseLoopVar

		t.Run(gitConfigScope.String(), func(t *testing.T) {
			t.Parallel()

			execSucceeds := func(scope scope.Scope, options ...string) ([]string, error) {
				if scope != gitConfigScope {
					t.Errorf("unexpected scope, expected: %s, got: %s", gitConfigScope, scope)
					t.Fail()
				}
				if options[0] != "--list" {
					t.Errorf("unexpected option, expected: %s, got: %s", "--list", options[0])
					t.Fail()
				}
				return lines, nil
			}

			settingsMap, _ := list(execSucceeds)(gitConfigScope)

			if !reflect.DeepEqual(settingsMap, expectedMap) {
				t.Errorf("expected: %s, got: %s", expectedMap, settingsMap)
				t.Fail()
			}
		})
	}
}

func TestListShouldReturnEmptyMapIfEmptyReturnFromGitConfigCommand(t *testing.T) {
	expectedMap := make(map[string]string, 0)
	execSucceedsEmpty := func(scope.Scope, ...string) ([]string, error) { return nil, nil }

	settingsMap, _ := list(execSucceedsEmpty)(scope.Global)

	if !reflect.DeepEqual(expectedMap, settingsMap) {
		t.Errorf("expected: %s, got: %s", expectedMap, settingsMap)
		t.Fail()
	}
}

func TestListShouldFailIfGitConfigCommandFails(t *testing.T) {
	expectedMap := make(map[string]string, 0)
	expectedErr := errors.New("failed to exec git config command")

	execFails := func(scope.Scope, ...string) ([]string, error) { return nil, expectedErr }

	settingsMap, err := list(execFails)(scope.Global)

	if !reflect.DeepEqual(expectedMap, settingsMap) {
		t.Errorf("expected: %s, got: %s", expectedMap, settingsMap)
		t.Fail()
	}

	if expectedErr != err {
		t.Errorf("expected err: %s, got err: %s", expectedErr, err)
		t.Fail()
	}
}
