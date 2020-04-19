package gitconfig

import (
	"reflect"
	"testing"
)

func TestFlag(t *testing.T) {
	t.Parallel()

	cases := []struct {
		scope        Scope
		expectedFlag string
	}{
		{Global, "--global"},
		{Local, "--local"},
	}

	for _, caseLoopVar := range cases {
		scope := caseLoopVar.scope
		expectedFlag := caseLoopVar.expectedFlag

		t.Run(scope.String(), func(t *testing.T) {
			t.Parallel()
			flag := scope.Flag()

			if !reflect.DeepEqual(expectedFlag, flag) {
				t.Errorf("expected: %s, got: %s", expectedFlag, scope)
				t.Fail()
			}
		})
	}
}
