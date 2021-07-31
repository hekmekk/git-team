package entity

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFromString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		candidateScope string
		expectedScope  Scope
	}{
		{"global", Global},
		{"repo-local", RepoLocal},
		{"unknown", Unknown},
		{"some other string", Unknown},
	}

	for _, caseLoopVar := range cases {
		candidate := caseLoopVar.candidateScope
		expectedScope := caseLoopVar.expectedScope

		t.Run(candidate, func(t *testing.T) {
			t.Parallel()
			scope := FromString(candidate)

			require.Equal(t, expectedScope, scope)
		})
	}
}
