package completion

import (
	"errors"
	"fmt"
	"testing"

	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
	"github.com/stretchr/testify/require"
)

type gitConfigReaderMock struct {
	getRegexp func(gitconfigscope.Scope, string) (map[string]string, error)
}

func (mock gitConfigReaderMock) Get(scope gitconfigscope.Scope, key string) (string, error) {
	return "", nil
}

func (mock gitConfigReaderMock) GetAll(scope gitconfigscope.Scope, key string) ([]string, error) {
	return []string{}, nil
}

func (mock gitConfigReaderMock) GetRegexp(scope gitconfigscope.Scope, pattern string) (map[string]string, error) {
	return mock.getRegexp(scope, pattern)
}

func (mock gitConfigReaderMock) List(scope gitconfigscope.Scope) (map[string]string, error) {
	return nil, nil
}

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

	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return map[string]string{
				"team.alias.alias1": "Mr. Noujz <noujz@mr.se>",
				"team.alias.alias2": "Mrs. Noujz <noujz@mrs.se>",
				"team.alias.alias3": "Mrs. Very Noujz <very-noujz@mrs.se>",
			}, nil
		},
	}

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
	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return map[string]string{}, nil
		},
	}

	expectedRemainingAliases := []string{}

	aliasShellCompletion := NewAliasShellCompletion(gitConfigReader)

	remainingAliases := aliasShellCompletion.Complete([]string{})

	require.Equal(t, expectedRemainingAliases, remainingAliases)
}

func TestCompletewhenLookingUpAssignmentsFails(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
			return map[string]string{}, errors.New("any kind of error")
		},
	}

	expectedRemainingAliases := []string{}

	aliasShellCompletion := NewAliasShellCompletion(gitConfigReader)

	remainingAliases := aliasShellCompletion.Complete([]string{})

	require.Equal(t, expectedRemainingAliases, remainingAliases)
}
