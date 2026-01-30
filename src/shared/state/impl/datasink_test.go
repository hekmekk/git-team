package stateimpl

import (
	mocks "github.com/hekmekk/git-team/v2/mocks/shared/gitconfig/interface"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"

	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	gitconfigerror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

func TestPersistSucceeds(t *testing.T) {
	gitConfigWriter := &mocks.Writer{}

	gitConfigWriter.
		On("UnsetAll", mock.Anything, mock.Anything).
		Return(nil)

	gitConfigWriter.
		On("Add", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	gitConfigWriter.
		On("ReplaceAll", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"}, "/previous/hooks/path")

	require.Nil(t, err)
}

func TestPersistSucceedsWhenTryingToRemoveNonExistingActiveCoauthorsFromGitConfig(t *testing.T) {
	gitConfigWriter := &mocks.Writer{}

	gitConfigWriter.
		On("UnsetAll", mock.Anything, mock.Anything).
		Return(gitconfigerror.ErrTryingToUnsetAnOptionWhichDoesNotExist)

	gitConfigWriter.
		On("Add", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	gitConfigWriter.
		On("ReplaceAll", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"}, "/previous/hooks/path")

	require.Nil(t, err)
}

func TestPersistFailsDueToAnotherUnsetAllFailure(t *testing.T) {
	gitConfigWriter := &mocks.Writer{}

	gitConfigWriter.
		On("UnsetAll", mock.Anything, mock.Anything).
		Return(gitconfigerror.ErrConfigFileCannotBeWritten)

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"}, "/previous/hooks/path")

	require.Error(t, err)
}

func TestPersistFailsDueToAddFailure(t *testing.T) {
	gitConfigWriter := &mocks.Writer{}

	gitConfigWriter.
		On("UnsetAll", mock.Anything, mock.Anything).
		Return(nil)

	gitConfigWriter.
		On("Add", mock.Anything, mock.Anything, mock.Anything).
		Return(gitconfigerror.ErrConfigFileCannotBeWritten)

	err := NewGitConfigDataSink(gitConfigWriter).PersistEnabled(activationscope.Global, []string{"CO-AUTHOR"}, "/previous/hooks/path")

	require.Error(t, err)
}

func TestPersistFailsDueReplaceAllFailure(t *testing.T) {
	gitConfigWriter := &mocks.Writer{}

	gitConfigWriter.
		On("UnsetAll", mock.Anything, mock.Anything).
		Return(nil)

	gitConfigWriter.
		On("Add", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	gitConfigWriter.
		On("ReplaceAll", mock.Anything, mock.Anything, mock.Anything).
		Return(gitconfigerror.ErrSectionOrKeyIsInvalid)

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

			gitConfigWriter := &mocks.Writer{}

			gitConfigWriter.
				On("UnsetAll", mock.Anything, mock.Anything).
				Return(nil)

			gitConfigWriter.
				On("Add", mock.Anything, mock.Anything, mock.Anything).
				Return(nil)

			gitConfigWriter.
				On("ReplaceAll", mock.Anything, mock.Anything, mock.Anything).
				Return(nil)

			err := NewGitConfigDataSink(gitConfigWriter).PersistDisabled(activationScope)

			require.Nil(t, err)

			gitConfigWriter.AssertNotCalled(t, "Add")

			gitConfigWriter.AssertCalled(t, "UnsetAll", expectedGitConfigScope, mock.Anything)
			gitConfigWriter.AssertCalled(t, "ReplaceAll", expectedGitConfigScope, mock.Anything, mock.Anything)
		})
	}

}
