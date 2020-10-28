package activationvalidatorimpl

import (
	"errors"
	"testing"

	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
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
