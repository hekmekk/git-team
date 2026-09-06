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
	"testing"

	scope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

func TestShouldReturnTheValue(t *testing.T) {
	expectedValue := "value"
	execSucceeds := func(scope.Scope, ...string) ([]string, error) { return []string{expectedValue}, nil }

	value, err := get(execSucceeds)(scope.Global, "key")

	if err != nil {
		t.Errorf("expected err: %s, got err: %s", "[no err]", err)
		t.Fail()
	}

	if expectedValue != value {
		t.Errorf("expected: %s, received: %s", expectedValue, value)
		t.Fail()
	}
}

func TestShouldReturnEmptyStringIfNoValueIsFound(t *testing.T) {
	execSucceedsEmpty := func(scope.Scope, ...string) ([]string, error) { return []string{}, nil }

	value, err := get(execSucceedsEmpty)(scope.Global, "key")

	if err != nil {
		t.Errorf("expected err: %s, got err: %s", "[no err]", err)
		t.Fail()
	}

	if "" != value {
		t.Errorf("expected: %s, received: %s", "", value)
		t.Fail()
	}
}

func TestShouldReturnTheErrorIfAnErrorOccurs(t *testing.T) {
	expectedErr := errors.New("git command failed")
	execFails := func(scope.Scope, ...string) ([]string, error) { return nil, expectedErr }

	value, err := get(execFails)(scope.Global, "key")

	if expectedErr != err {
		t.Errorf("expected err: %s, got err: %s", expectedErr, err)
		t.Fail()
	}

	if "" != value {
		t.Errorf("expected: %s, received: %s", "", value)
		t.Fail()
	}
}
