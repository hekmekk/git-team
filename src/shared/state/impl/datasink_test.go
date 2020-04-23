package stateimpl

import (
	"errors"
	"testing"

	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

type gitConfigWriterMock struct {
	unsetAll   func(gitconfigscope.Scope, string) error
	replaceAll func(gitconfigscope.Scope, string, string) error
	add        func(gitconfigscope.Scope, string, string) error
}

func (mock gitConfigWriterMock) UnsetAll(scope gitconfigscope.Scope, key string) error {
	return mock.unsetAll(scope, key)
}

func (mock gitConfigWriterMock) ReplaceAll(scope gitconfigscope.Scope, key string, value string) error {
	return mock.replaceAll(scope, key, value)
}

func (mock gitConfigWriterMock) Add(scope gitconfigscope.Scope, key string, value string) error {
	return mock.add(scope, key, value)
}

func TestPersistSucceeds(t *testing.T) {
	gitConfigWriter := &gitConfigWriterMock{
		unsetAll:   func(gitconfigscope.Scope, string) error { return nil },
		add:        func(gitconfigscope.Scope, string, string) error { return nil },
		replaceAll: func(gitconfigscope.Scope, string, string) error { return nil },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"})

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestPersistSucceedsWhenTryingToRemoveNonExistingActiveCoauthorsFromGitConfig(t *testing.T) {
	gitConfigWriter := &gitConfigWriterMock{
		unsetAll:   func(gitconfigscope.Scope, string) error { return errors.New("exec status 5") },
		add:        func(gitconfigscope.Scope, string, string) error { return nil },
		replaceAll: func(gitconfigscope.Scope, string, string) error { return nil },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"})

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestPersistFailsDueToAnotherUnsetAllFailure(t *testing.T) {
	expectedErr := errors.New("unset-all failure")

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(gitconfigscope.Scope, string) error { return expectedErr },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"})

	if expectedErr != err {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}
}

func TestPersistFailsDueToAddFailure(t *testing.T) {
	expectedErr := errors.New("add failure")

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(gitconfigscope.Scope, string) error { return nil },
		add:      func(gitconfigscope.Scope, string, string) error { return expectedErr },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"})

	if expectedErr != err {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}
}

func TestPersistFailsDueReplaceAllFailure(t *testing.T) {
	expectedErr := errors.New("replace-all failure")

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll:   func(gitconfigscope.Scope, string) error { return nil },
		add:        func(gitconfigscope.Scope, string, string) error { return nil },
		replaceAll: func(gitconfigscope.Scope, string, string) error { return expectedErr },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistDisabled(activationscope.Global)

	if expectedErr != err {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}
}

func TestPersistPassesThroughTheCorrectScope(t *testing.T) {
	t.Parallel()

	cases := []struct {
		activationScope activationscope.ActivationScope
		gitConfigScope  gitconfigscope.Scope
	}{
		{activationscope.Global, gitconfigscope.Global},
		{activationscope.RepoLocal, gitconfigscope.Local},
	}

	for _, caseLoopVar := range cases {
		activationScope := caseLoopVar.activationScope
		expectedGitConfigScope := caseLoopVar.gitConfigScope

		t.Run(activationScope.String(), func(t *testing.T) {
			t.Parallel()

			gitConfigWriter := &gitConfigWriterMock{
				unsetAll: func(actualScope gitconfigscope.Scope, _ string) error {
					if actualScope != expectedGitConfigScope {
						t.Errorf("unsetAll - expected: %s, got: %s", expectedGitConfigScope, actualScope)
						t.Fail()
					}
					return nil
				},
				add: func(actualScope gitconfigscope.Scope, _ string, _ string) error {
					if actualScope != expectedGitConfigScope {
						t.Errorf("add - expected: %s, got: %s", expectedGitConfigScope, actualScope)
						t.Fail()
					}
					return nil
				},
				replaceAll: func(actualScope gitconfigscope.Scope, _ string, _ string) error {
					if actualScope != expectedGitConfigScope {
						t.Errorf("replaceAll - expected: %s, got: %s", expectedGitConfigScope, actualScope)
						t.Fail()
					}
					return nil
				},
			}

			err := NewGitConfigDataSink(gitConfigWriter).PersistDisabled(activationScope)

			if err != nil {
				t.Error(err)
				t.Fail()
			}
		})
	}

}
