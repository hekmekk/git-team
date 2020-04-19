package datasink

import (
	"errors"
	"reflect"
	"testing"

	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
)

type gitConfigWriterMock struct {
	replaceAll func(key string, value string) error
}

func (mock gitConfigWriterMock) ReplaceAll(key string, value string) error {
	return mock.replaceAll(key, value)
}

func (mock gitConfigWriterMock) UnsetAll(key string) error {
	return nil
}

func TestSetActivationScopeSucceeds(t *testing.T) {
	gitConfigWriter := gitConfigWriterMock{
		replaceAll: func(key string, value string) error {
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
		replaceAll: func(key string, value string) error {
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
