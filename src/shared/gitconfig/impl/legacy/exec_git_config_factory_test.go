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

func TestShouldExecuteGitConfigWithTheExpectedCommandLineArguments(t *testing.T) {
	t.Parallel()

	cases := []struct {
		scope             scope.Scope
		expectedScopeFlag string
	}{
		{scope.Global, "--global"},
		{scope.Local, "--local"},
	}

	for _, caseLoopVar := range cases {
		scope := caseLoopVar.scope
		expectedScopeFlag := caseLoopVar.expectedScopeFlag

		t.Run(scope.String(), func(t *testing.T) {
			providedOptions := []string{"--one", "--two"}
			expectedOptions := append([]string{expectedScopeFlag}, providedOptions...)

			executor := func(args ...string) (string, error) {
				if !reflect.DeepEqual(expectedOptions, args) {
					t.Errorf("expected: %s, received: %s", expectedOptions, args)
					t.Fail()
				}

				return "", nil
			}

			execGitConfigFactory(executor)(scope, providedOptions...)
		})
	}
}

func TestShouldReturnTheTwoLines(t *testing.T) {
	expectedLines := []string{"line1", "line2"}
	out := fmt.Sprintf("%s\n%s\n", expectedLines[0], expectedLines[1])

	executor := func(args ...string) (string, error) {
		return out, nil
	}

	lines, err := execGitConfigFactory(executor)(scope.Global, "")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedLines, lines) {
		t.Errorf("expected %s, received: %s", expectedLines, lines)
		t.Fail()
	}
}

func TestShouldReturnNoLines(t *testing.T) {
	expectedLines := []string{}

	executor := func(args ...string) (string, error) {
		return "", nil
	}

	lines, err := execGitConfigFactory(executor)(scope.Global, "")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedLines, lines) {
		t.Errorf("expected %s, received: %s", expectedLines, lines)
		t.Fail()
	}
}

func TestShouldReturnAnError(t *testing.T) {
	expectedLines := []string{}

	executor := func(args ...string) (string, error) {
		return "", errors.New("executing the git config command failed")
	}

	lines, err := execGitConfigFactory(executor)(scope.Global, "")
	if err == nil {
		t.Error("expected an error")
		t.Fail()
	}

	if lines != nil {
		t.Errorf("expected %s, received: %s", expectedLines, lines)
		t.Fail()
	}
}
