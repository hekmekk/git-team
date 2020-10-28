package stateimpl

import (
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
			return "disabled", nil
		},
		getAll: func(scope gitconfigscope.Scope, key string) ([]string, error) {
			return []string{}, nil
		},
	}

	state, err := NewGitConfigDataSource(gitConfigReader).Query(activationscope.Global)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, state) {
		t.Errorf("expected: %s, got: %s", expectedState, state)
		t.Fail()
	}
}

func TestQueryEnabled(t *testing.T) {
	activeCoauthors := []string{"Mr. Noujz <noujz@mr.se>"}
	expectedState := state.State{Status: "enabled", Coauthors: activeCoauthors}

	gitConfigReader := &gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			return "enabled", nil
		},
		getAll: func(scope gitconfigscope.Scope, key string) ([]string, error) {
			return activeCoauthors, nil
		},
	}

	state, err := NewGitConfigDataSource(gitConfigReader).Query(activationscope.Global)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, state) {
		t.Errorf("expected: %s, got: %s", expectedState, state)
		t.Fail()
	}
}

func TestQueryDisabledWhenStatusUnset(t *testing.T) {
	expectedState := state.State{Status: "disabled", Coauthors: []string{}}

	gitConfigReader := &gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			return "", nil
		},
		getAll: func(scope gitconfigscope.Scope, key string) ([]string, error) {
			return []string{}, nil
		},
	}

	state, err := NewGitConfigDataSource(gitConfigReader).Query(activationscope.Global)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, state) {
		t.Errorf("expected: %s, got: %s", expectedState, state)
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
					return "enabled", nil
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
