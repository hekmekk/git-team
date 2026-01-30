package stateimpl

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	gitconfigerror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

type gitConfigWriterMock struct {
	add        func(gitconfigscope.Scope, string, string) error
	unsetAll   func(gitconfigscope.Scope, string) error
	replaceAll func(gitconfigscope.Scope, string, string) error
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
		add:        func(_ gitconfigscope.Scope, _ string, _ string) error { return nil },
		unsetAll:   func(_ gitconfigscope.Scope, _ string) error { return nil },
		replaceAll: func(_ gitconfigscope.Scope, _ string, _ string) error { return nil },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"}, "/previous/hooks/path")

	require.Nil(t, err)
}

func TestPersistSucceedsWhenTryingToRemoveNonExistingActiveCoauthorsFromGitConfig(t *testing.T) {
	gitConfigWriter := &gitConfigWriterMock{
		add: func(_ gitconfigscope.Scope, _ string, _ string) error { return nil },
		unsetAll: func(_ gitconfigscope.Scope, _ string) error {
			return gitconfigerror.ErrTryingToUnsetAnOptionWhichDoesNotExist
		},
		replaceAll: func(_ gitconfigscope.Scope, _ string, _ string) error { return nil },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"}, "/previous/hooks/path")

	require.Nil(t, err)
}

func TestPersistFailsDueToAnotherUnsetAllFailure(t *testing.T) {
	gitConfigWriter := &gitConfigWriterMock{
		add: func(_ gitconfigscope.Scope, _ string, _ string) error { return nil },
		unsetAll: func(_ gitconfigscope.Scope, _ string) error {
			return gitconfigerror.ErrConfigFileCannotBeWritten
		},
		replaceAll: func(_ gitconfigscope.Scope, _ string, _ string) error { return nil },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"}, "/previous/hooks/path")

	require.Error(t, err)
}

func TestPersistFailsDueToAddFailure(t *testing.T) {
	gitConfigWriter := &gitConfigWriterMock{
		add: func(_ gitconfigscope.Scope, _ string, _ string) error {
			return gitconfigerror.ErrConfigFileCannotBeWritten
		},
		unsetAll:   func(_ gitconfigscope.Scope, _ string) error { return nil },
		replaceAll: func(_ gitconfigscope.Scope, _ string, _ string) error { return nil },
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"}, "/previous/hooks/path")

	require.Error(t, err)
}

func TestPersistFailsDueReplaceAllFailure(t *testing.T) {
	gitConfigWriter := &gitConfigWriterMock{
		add:      func(_ gitconfigscope.Scope, _ string, _ string) error { return nil },
		unsetAll: func(_ gitconfigscope.Scope, _ string) error { return nil },
		replaceAll: func(_ gitconfigscope.Scope, _ string, _ string) error {
			return gitconfigerror.ErrSectionOrKeyIsInvalid
		},
	}

	err := NewGitConfigDataSink(gitConfigWriter).PersistDisabled(activationscope.Global)

	require.Error(t, err)
}

func TestPersistPassesThroughTheCorrectScope(t *testing.T) {
	t.Parallel()

	cases := []struct {
		activationScope activationscope.Scope
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
				add: func(s gitconfigscope.Scope, _ string, _ string) error {
					return errors.New("should not have been called")
				},
				unsetAll: func(actualScope gitconfigscope.Scope, _ string) error {
					require.Equal(t, expectedGitConfigScope, actualScope)
					return nil
				},
				replaceAll: func(actualScope gitconfigscope.Scope, _ string, _ string) error {
					require.Equal(t, expectedGitConfigScope, actualScope)
					return nil
				},
			}

			err := NewGitConfigDataSink(gitConfigWriter).PersistDisabled(activationScope)

			require.Nil(t, err)
		})
	}

}
