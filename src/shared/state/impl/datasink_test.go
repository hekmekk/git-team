package stateimpl

import (
	"errors"
	"testing"

	scope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

type gitConfigWriterMock struct {
	unsetAll   func(scope.Scope, string) error
	replaceAll func(scope.Scope, string, string) error
	add        func(scope.Scope, string, string) error
}

func (mock gitConfigWriterMock) UnsetAll(scope scope.Scope, key string) error {
	return mock.unsetAll(scope, key)
}

func (mock gitConfigWriterMock) ReplaceAll(scope scope.Scope, key string, value string) error {
	return mock.replaceAll(scope, key, value)
}

func (mock gitConfigWriterMock) Add(scope scope.Scope, key string, value string) error {
	return mock.add(scope, key, value)
}

func TestPersistSucceeds(t *testing.T) {
	gitConfigWriter := &gitConfigWriterMock{
		unsetAll:   func(scope.Scope, string) error { return nil },
		add:        func(scope.Scope, string, string) error { return nil },
		replaceAll: func(scope.Scope, string, string) error { return nil },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(scope.Global, []string{"CO-AUTHOR"})

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestPersistSucceedsWhenTryingToRemoveNonExistingActiveCoauthorsFromGitConfig(t *testing.T) {
	gitConfigWriter := &gitConfigWriterMock{
		unsetAll:   func(scope.Scope, string) error { return errors.New("exec status 5") },
		add:        func(scope.Scope, string, string) error { return nil },
		replaceAll: func(scope.Scope, string, string) error { return nil },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(scope.Global, []string{"CO-AUTHOR"})

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestPersistFailsDueToAnotherUnsetAllFailure(t *testing.T) {
	expectedErr := errors.New("unset-all failure")

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(scope.Scope, string) error { return expectedErr },
		// add:        func(scope.Scope, string, string) error { return nil },
		// replaceAll: func(scope.Scope, string, string) error { return nil },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(scope.Global, []string{"CO-AUTHOR"})

	if expectedErr != err {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}
}

func TestPersistFailsDueToAddFailure(t *testing.T) {
	expectedErr := errors.New("add failure")

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(scope.Scope, string) error { return nil },
		add:      func(scope.Scope, string, string) error { return expectedErr },
		// replaceAll: func(scope.Scope, string, string) error { return nil },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(scope.Global, []string{"CO-AUTHOR"})

	if expectedErr != err {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}
}

func TestPersistFailsDueReplaceAllFailure(t *testing.T) {
	expectedErr := errors.New("replace-all failure")

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll:   func(scope.Scope, string) error { return nil },
		add:        func(scope.Scope, string, string) error { return nil },
		replaceAll: func(scope.Scope, string, string) error { return expectedErr },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistDisabled(scope.Global)

	if expectedErr != err {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}
}

func TestPersistPassesThroughTheCorrectScope(t *testing.T) {
	t.Parallel()

	for _, caseLoopVar := range []scope.Scope{scope.Global, scope.Local} {
		theScope := caseLoopVar

		t.Run(theScope.String(), func(t *testing.T) {
			t.Parallel()

			gitConfigWriter := &gitConfigWriterMock{
				unsetAll: func(passThroughScope scope.Scope, _ string) error {
					if passThroughScope != theScope {
						t.Errorf("unsetAll - expected: %s, got: %s", theScope, passThroughScope)
						t.Fail()
					}
					return nil
				},
				add: func(passThroughScope scope.Scope, _ string, _ string) error {
					if passThroughScope != theScope {
						t.Errorf("add - expected: %s, got: %s", theScope, passThroughScope)
						t.Fail()
					}
					return nil
				},
				replaceAll: func(passThroughScope scope.Scope, _ string, _ string) error {
					if passThroughScope != theScope {
						t.Errorf("replaceAll - expected: %s, got: %s", theScope, passThroughScope)
						t.Fail()
					}
					return nil
				},
			}

			err := NewGitConfigDataSink(gitConfigWriter).PersistDisabled(theScope)

			if err != nil {
				t.Error(err)
				t.Fail()
			}
		})
	}

}
