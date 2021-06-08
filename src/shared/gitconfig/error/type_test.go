package gitconfigerror

import (
	"errors"
	"reflect"
	"testing"
)

func TestConstructor(t *testing.T) {
	cases := []struct {
		description          string
		originalErr          error
		expectedGitConfigErr error
	}{
		{description: "no error", originalErr: nil, expectedGitConfigErr: nil},
		{description: "unknown error", originalErr: errors.New("unknown"), expectedGitConfigErr: &UnknownErr{message: "unknown gitconfig error: unknown"}},
		{description: ErrSectionOrKeyIsInvalid.Error(), originalErr: errors.New("exit status 1"), expectedGitConfigErr: ErrSectionOrKeyIsInvalid},
		{description: ErrNoSectionOrNameProvided.Error(), originalErr: errors.New("exit status 2"), expectedGitConfigErr: ErrNoSectionOrNameProvided},
		{description: ErrConfigFileIsInvalid.Error(), originalErr: errors.New("exit status 3"), expectedGitConfigErr: ErrConfigFileIsInvalid},
		{description: ErrConfigFileCannotBeWritten.Error(), originalErr: errors.New("exit status 4"), expectedGitConfigErr: ErrConfigFileCannotBeWritten},
		{description: ErrTryingToUnsetAnOptionWhichDoesNotExist.Error(), originalErr: errors.New("exit status 5"), expectedGitConfigErr: ErrTryingToUnsetAnOptionWhichDoesNotExist},
		{description: ErrTryingToUseAnInvalidRegexp.Error(), originalErr: errors.New("exit status 6"), expectedGitConfigErr: ErrTryingToUseAnInvalidRegexp},
	}

	for _, caseLoopVar := range cases {
		originalErr := caseLoopVar.originalErr
		expectedGitConfigErr := caseLoopVar.expectedGitConfigErr
		description := caseLoopVar.description

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			actualGitConfigErr := New(originalErr)

			if !reflect.DeepEqual(expectedGitConfigErr, actualGitConfigErr) {
				t.Errorf("expected: %s, got: %s", expectedGitConfigErr, actualGitConfigErr)
				t.Fail()
			}
		})
	}

}
