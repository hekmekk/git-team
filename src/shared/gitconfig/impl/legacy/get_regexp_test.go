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

func TestShouldReturnTheKeyValueMap(t *testing.T) {
	key1 := "key1"
	val1 := "space separated value 1"
	key2 := "key2"
	val2 := "space separated value 2"

	expectedMap := make(map[string]string)
	expectedMap[key1] = val1
	expectedMap[key2] = val2

	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("%s %s", key1, val1))
	lines = append(lines, fmt.Sprintf("%s %s", key2, val2))

	execSucceeds := func(scope.Scope, ...string) ([]string, error) { return lines, nil }

	aliasMap, _ := getRegexp(execSucceeds)(scope.Global, "pattern")

	if !reflect.DeepEqual(aliasMap, expectedMap) {
		t.Errorf("expected: %s, got: %s", expectedMap, aliasMap)
		t.Fail()
	}
}

func TestShouldReturnEmptyMapIfEmptyReturnFromGitConfigCommand(t *testing.T) {
	expectedMap := make(map[string]string, 0)
	execSucceedsEmpty := func(scope.Scope, ...string) ([]string, error) { return nil, nil }

	aliasMap, _ := getRegexp(execSucceedsEmpty)(scope.Global, "pattern")

	if !reflect.DeepEqual(expectedMap, aliasMap) {
		t.Errorf("expected: %s, got: %s", expectedMap, aliasMap)
		t.Fail()
	}
}

func TestShouldFailIfGitConfigCommandFails(t *testing.T) {
	expectedMap := make(map[string]string, 0)
	expectedErr := errors.New("failed to exec git config command")

	execFails := func(scope.Scope, ...string) ([]string, error) { return nil, expectedErr }

	aliasMap, err := getRegexp(execFails)(scope.Global, "pattern")

	if !reflect.DeepEqual(expectedMap, aliasMap) {
		t.Errorf("expected: %s, got: %s", expectedMap, aliasMap)
		t.Fail()
	}

	if expectedErr != err {
		t.Errorf("expected err: %s, got err: %s", expectedErr, err)
		t.Fail()
	}
}
