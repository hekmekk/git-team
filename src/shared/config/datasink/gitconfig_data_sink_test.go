package datasink

import (
	"errors"
	"reflect"
	"testing"

	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

type gitConfigWriterMock struct {
	replaceAll func(scope gitconfigscope.Scope, key string, value string) error
}

func (mock gitConfigWriterMock) Add(scope gitconfigscope.Scope, key string, value string) error {
	return nil
}

func (mock gitConfigWriterMock) ReplaceAll(scope gitconfigscope.Scope, key string, value string) error {
	return mock.replaceAll(scope, key, value)
}

func (mock gitConfigWriterMock) UnsetAll(scope gitconfigscope.Scope, key string) error {
	return nil
}

func TestSetActivationScopeSucceeds(t *testing.T) {
	gitConfigWriter := gitConfigWriterMock{
		replaceAll: func(scope gitconfigscope.Scope, key string, value string) error {
			if scope != gitconfigscope.Global {
				return errors.New("wrong scope")
			}
			if key != "team.config.activation-scope" {
				return errors.New("wrong key")
			}
			return nil
		},
	}

	err := NewGitconfigDataSink(gitConfigWriter).SetActivationScope(activationscope.RepoLocal)

	if err != nil {
		t.Errorf("expected: no error, received: '%s'", err)
		t.Fail()
	}
}

func TestSetActivationScopeFails(t *testing.T) {
	expectedErr := errors.New("gitconfig failed")

	gitConfigWriter := gitConfigWriterMock{
		replaceAll: func(scope gitconfigscope.Scope, key string, value string) error {
			if scope != gitconfigscope.Global {
				return errors.New("wrong scope")
			}
			if key != "team.config.activation-scope" {
				return errors.New("wrong key")
			}
			return errors.New("gitconfig failed")
		},
	}

	err := NewGitconfigDataSink(gitConfigWriter).SetActivationScope(activationscope.RepoLocal)

	if !reflect.DeepEqual(expectedErr, err) {
		t.Errorf("expected: '%s', received: '%s'", expectedErr, err)
		t.Fail()
	}
}
