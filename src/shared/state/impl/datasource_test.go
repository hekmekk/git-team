package stateimpl

import (
	"fmt"
	gitconfigerror "github.com/hekmekk/git-team/src/shared/gitconfig/error"
	"reflect"
	"testing"

	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
	state "github.com/hekmekk/git-team/src/shared/state/entity"
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
