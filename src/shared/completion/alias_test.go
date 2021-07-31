package completion

import (
	"errors"
	"fmt"
	mocks "github.com/hekmekk/git-team/mocks/shared/gitconfig/interface"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestComplete(t *testing.T) {
	t.Parallel()

	cases := []struct {
		selectedAliases          []string
		expectedRemainingAliases []string
	}{
		{[]string{}, []string{"alias1", "alias2", "alias3"}},
		{[]string{"alias1"}, []string{"alias2", "alias3"}},
		{[]string{"alias1", "alias2"}, []string{"alias3"}},
		{[]string{"alias1", "alias2", "alias3"}, []string{}},
	}

	gitConfigReader := &mocks.Reader{}

	gitConfigReader.On("GetRegexp", mock.Anything, mock.Anything).Return(map[string]string{
		"team.alias.alias1": "Mr. Noujz <noujz@mr.se>",
		"team.alias.alias2": "Mrs. Noujz <noujz@mrs.se>",
		"team.alias.alias3": "Mrs. Very Noujz <very-noujz@mrs.se>",
	}, nil)

	aliasShellCompletion := NewAliasShellCompletion(gitConfigReader)

	for _, caseLoopVar := range cases {
		selectedAliases := caseLoopVar.selectedAliases
		expectedRemainingAliases := caseLoopVar.expectedRemainingAliases

		t.Run(fmt.Sprintf("selectedAliases: %s expectedRemainingAliases: %s", selectedAliases, expectedRemainingAliases), func(t *testing.T) {
			t.Parallel()

			remainingAliases := aliasShellCompletion.Complete(selectedAliases)

			require.Equal(t, expectedRemainingAliases, remainingAliases)
		})
	}
}

func TestCompleteWhenNoAssignmentsExists(t *testing.T) {
	gitConfigReader := &mocks.Reader{}

	gitConfigReader.On("GetRegexp", mock.Anything, mock.Anything).Return(map[string]string{}, nil)

	expectedRemainingAliases := []string{}

	aliasShellCompletion := NewAliasShellCompletion(gitConfigReader)

	remainingAliases := aliasShellCompletion.Complete([]string{})

	require.Equal(t, expectedRemainingAliases, remainingAliases)
}

func TestCompletewhenLookingUpAssignmentsFails(t *testing.T) {
	gitConfigReader := &mocks.Reader{}

	gitConfigReader.On("GetRegexp", mock.Anything, mock.Anything).Return(map[string]string{}, errors.New("any kind of error"))

	expectedRemainingAliases := []string{}

	aliasShellCompletion := NewAliasShellCompletion(gitConfigReader)

	remainingAliases := aliasShellCompletion.Complete([]string{})

	require.Equal(t, expectedRemainingAliases, remainingAliases)
}
