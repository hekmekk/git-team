package datasink

import (
	"errors"
	"reflect"
	"testing"

	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
)

type gitconfigRawWriterMock struct {
	set func(key string, value string) error
}

func (mock gitconfigRawWriterMock) Set(key string, value string) error {
	return mock.set(key, value)
}

func (mock gitconfigRawWriterMock) Unset(key string) error {
	return nil
}

func TestSetActivationScopeSucceeds(t *testing.T) {
	rawWriter := gitconfigRawWriterMock{
		set: func(key string, value string) error {
			if key != "team.config.activation-scope" {
				return errors.New("wrong key")
			}
			return nil
		},
	}

	err := newGitconfigDataSink(rawWriter).SetActivationScope(activationscope.RepoLocal)

	if err != nil {
		t.Errorf("expected: no error, received: '%s'", err)
		t.Fail()
	}
}

func TestSetActivationScopeFails(t *testing.T) {
	expectedErr := errors.New("gitconfig failed")

	rawWriter := gitconfigRawWriterMock{
		set: func(key string, value string) error {
			if key != "team.config.activation-scope" {
				return errors.New("wrong key")
			}
			return errors.New("gitconfig failed")
		},
	}

	err := newGitconfigDataSink(rawWriter).SetActivationScope(activationscope.RepoLocal)

	if !reflect.DeepEqual(expectedErr, err) {
		t.Errorf("expected: '%s', received: '%s'", expectedErr, err)
		t.Fail()
	}
}
