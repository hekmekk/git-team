package stateimpl

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
	gitconfigerror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
	"reflect"
	"testing"

	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
	state "github.com/hekmekk/git-team/v2/src/shared/state/entity"
)

type gitConfigReaderMock struct {
	get    func(gitconfigscope.Scope, string) (string, error)
	getAll func(gitconfigscope.Scope, string) ([]string, error)
}

func (mock gitConfigReaderMock) Get(scope gitconfigscope.Scope, key string) (string, error) {
	return mock.get(scope, key)
}

func (mock gitConfigReaderMock) GetAll(scope gitconfigscope.Scope, key string) ([]string, error) {
	return mock.getAll(scope, key)
}

func (mock gitConfigReaderMock) GetRegexp(scope gitconfigscope.Scope, pattern string) (map[string]string, error) {
	return nil, nil
}

func (mock gitConfigReaderMock) List(scope gitconfigscope.Scope) (map[string]string, error) {
	return nil, nil
}

func TestQueryDisabled(t *testing.T) {
	expectedState := state.State{Status: "disabled", Coauthors: []string{}}

	gitConfigReader := &gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			if "team.state.status" == key {
				return "disabled", nil
			}

			if "team.state.previous-hooks-path" == key {
				return "", nil
			}

			return "", nil
		},
		getAll: func(scope gitconfigscope.Scope, key string) ([]string, error) {
			return []string{}, nil
		},
	}

	appState, err := NewGitConfigDataSource(gitConfigReader).Query(activationscope.Global)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, appState) {
		t.Errorf("expected: %s, got: %s", expectedState, appState)
		t.Fail()
	}
}

func TestQueryEnabled(t *testing.T) {
	activeCoauthors := []string{"Mr. Noujz <noujz@mr.se>"}
	expectedState := state.State{Status: "enabled", Coauthors: activeCoauthors, PreviousHooksPath: "/path/to/previous/hooks"}

	gitConfigReader := &gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			if "team.state.status" == key {
				return "enabled", nil
			}

			if "team.state.previous-hooks-path" == key {
				return "/path/to/previous/hooks", nil
			}

			return "", nil
		},
		getAll: func(scope gitconfigscope.Scope, key string) ([]string, error) {
			return activeCoauthors, nil
		},
	}

	appState, err := NewGitConfigDataSource(gitConfigReader).Query(activationscope.Global)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, appState) {
		t.Errorf("expected: %s, got: %s", expectedState, appState)
		t.Fail()
	}
}

func TestQueryDisabledWhenStatusUnset(t *testing.T) {
	expectedState := state.State{Status: "disabled", Coauthors: []string{}}

	gitConfigReader := &gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			if "team.state.status" == key {
				return "", nil
			}

			if "team.state.previous-hooks-path" == key {
				return "", nil
			}

			return "", nil
		},
		getAll: func(scope gitconfigscope.Scope, key string) ([]string, error) {
			return []string{}, nil
		},
	}

	appState, err := NewGitConfigDataSource(gitConfigReader).Query(activationscope.Global)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, appState) {
		t.Errorf("expected: %s, got: %s", expectedState, appState)
		t.Fail()
	}
}

func TestQueryTranslatesActivationScopeToGitconfigScopeCorrectly(t *testing.T) {
	t.Parallel()

	properties := []struct {
		activationScope activationscope.Scope
		gitConfigScope  gitconfigscope.Scope
	}{
		{activationscope.Global, gitconfigscope.Global},
		{activationscope.RepoLocal, gitconfigscope.Local},
	}

	for _, caseLoopVar := range properties {
		activationScope := caseLoopVar.activationScope
		gitConfigScope := caseLoopVar.gitConfigScope

		t.Run(activationScope.String(), func(t *testing.T) {
			t.Parallel()

			gitConfigReader := &gitConfigReaderMock{
				get: func(scope gitconfigscope.Scope, key string) (string, error) {
					if scope != gitConfigScope {
						t.Errorf("wrong scope, expected: %s, got: %s", gitConfigScope, scope)
						t.Fail()
					}
					if "team.state.status" == key {
						return "enabled", nil
					}

					if "team.state.previous-hooks-path" == key {
						return "", nil
					}

					return "", nil
				},
				getAll: func(scope gitconfigscope.Scope, key string) ([]string, error) {
					if scope != gitConfigScope {
						t.Errorf("wrong scope, expected: %s, got: %s", gitConfigScope, scope)
						t.Fail()
					}
					return []string{"Mr. Noujz <noujz@mr.se>"}, nil
				},
			}

			NewGitConfigDataSource(gitConfigReader).Query(activationScope)
		})
	}
}

func TestFailWhenFailingToRetrieveActiveCoAuthors(t *testing.T) {
	expectedState := state.State{}

	gitConfigReader := &gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			if "team.state.status" == key {
				return "enabled", nil
			}

			if "team.state.previous-hooks-path" == key {
				return "", nil
			}

			return "", nil
		},
		getAll: func(scope gitconfigscope.Scope, key string) ([]string, error) {
			return []string{}, gitconfigerror.ErrSectionOrKeyIsInvalid
		},
	}

	appState, err := NewGitConfigDataSource(gitConfigReader).Query(activationscope.Global)

	expectedErr := fmt.Errorf("no active co-authors found: %s", gitconfigerror.ErrSectionOrKeyIsInvalid)

	if !reflect.DeepEqual(expectedErr, err) {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, appState) {
		t.Errorf("expected: %s, got: %s", expectedState, appState)
		t.Fail()
	}
}

func TestEmptyHooksPathOnSpecificGitConfigErr(t *testing.T) {
	expectedState := state.State{Status: "enabled", Coauthors: []string{}}

	gitConfigReader := &gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			if "team.state.status" == key {
				return "enabled", nil
			}

			if "team.state.previous-hooks-path" == key {
				return "", gitconfigerror.ErrSectionOrKeyIsInvalid
			}

			return "", nil
		},
		getAll: func(scope gitconfigscope.Scope, key string) ([]string, error) {
			return []string{}, nil
		},
	}

	appState, err := NewGitConfigDataSource(gitConfigReader).Query(activationscope.Global)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, appState) {
		t.Errorf("expected: %s, got: %s", expectedState, appState)
		t.Fail()
	}
}

func TestFailWhenFailingToRetrievePreviousHooksPath(t *testing.T) {
	expectedState := state.State{}

	gitConfigReader := &gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			if "team.state.status" == key {
				return "enabled", nil
			}

			if "team.state.previous-hooks-path" == key {
				return "", gitconfigerror.ErrConfigFileIsInvalid
			}

			return "", nil
		},
		getAll: func(scope gitconfigscope.Scope, key string) ([]string, error) {
			return []string{}, nil
		},
	}

	appState, err := NewGitConfigDataSource(gitConfigReader).Query(activationscope.Global)

	expectedErr := fmt.Errorf("failed to get previous hooks path: %s", gitconfigerror.ErrConfigFileIsInvalid)

	if !reflect.DeepEqual(expectedErr, err) {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, appState) {
		t.Errorf("expected: %s, got: %s", expectedState, appState)
		t.Fail()
	}
}
