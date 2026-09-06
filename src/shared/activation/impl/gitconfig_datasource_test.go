package activationvalidatorimpl

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

	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

type gitConfigReaderMock struct {
	list func(gitconfigscope.Scope) (map[string]string, error)
}

func (mock gitConfigReaderMock) Get(scope gitconfigscope.Scope, key string) (string, error) {
	return "", nil
}

func (mock gitConfigReaderMock) GetAll(scope gitconfigscope.Scope, key string) ([]string, error) {
	return []string{}, nil
}

func (mock gitConfigReaderMock) GetRegexp(scope gitconfigscope.Scope, pattern string) (map[string]string, error) {
	return nil, nil
}

func (mock gitConfigReaderMock) List(scope gitconfigscope.Scope) (map[string]string, error) {
	return mock.list(scope)
}

func TestIsInsideAGitRepositoryTrue(t *testing.T) {
	gitConfigReader := gitConfigReaderMock{
		list: func(scope gitconfigscope.Scope) (map[string]string, error) {
			return make(map[string]string, 0), nil
		},
	}

	isInsideAGitRepository := NewGitConfigDataSource(gitConfigReader).IsInsideAGitRepository()

	if isInsideAGitRepository != true {
		t.Errorf("expected: %t, received %t", true, isInsideAGitRepository)
		t.Fail()
	}
}

func TestIsInsideAGitRepositoryFalse(t *testing.T) {
	gitConfigReader := gitConfigReaderMock{
		list: func(scope gitconfigscope.Scope) (map[string]string, error) {
			return make(map[string]string, 0), errors.New("some error")
		},
	}

	isInsideAGitRepository := NewGitConfigDataSource(gitConfigReader).IsInsideAGitRepository()

	if isInsideAGitRepository != false {
		t.Errorf("expected: %t, received %t", false, isInsideAGitRepository)
		t.Fail()
	}
}
