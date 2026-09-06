package gitconfigerror

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"errors"
	"reflect"
	"testing"
)

func TestConstructor(t *testing.T) {
	cases := []struct {
		description          string
		originalErr          error
		gitCmd               string
		gitOutput            string
		expectedGitConfigErr error
	}{
		{description: "no error", originalErr: nil, expectedGitConfigErr: nil},
		{description: "unknown error", originalErr: errors.New("unknown"), gitCmd: "the actual git-config command", gitOutput: "git-config-err line 1\ngit-config-err line 2", expectedGitConfigErr: &UnknownErr{message: "unknown gitconfig error: unknown\n\ngit-config command:\nthe actual git-config command\n\ngit-config output:\ngit-config-err line 1\ngit-config-err line 2"}},
		{description: ErrSectionOrKeyIsInvalid.Error(), originalErr: errors.New("exit status 1"), expectedGitConfigErr: ErrSectionOrKeyIsInvalid},
		{description: ErrNoSectionOrNameProvided.Error(), originalErr: errors.New("exit status 2"), expectedGitConfigErr: ErrNoSectionOrNameProvided},
		{description: ErrConfigFileIsInvalid.Error(), originalErr: errors.New("exit status 3"), expectedGitConfigErr: ErrConfigFileIsInvalid},
		{description: ErrConfigFileCannotBeWritten.Error(), originalErr: errors.New("exit status 4"), expectedGitConfigErr: ErrConfigFileCannotBeWritten},
		{description: ErrTryingToUnsetAnOptionWhichDoesNotExist.Error(), originalErr: errors.New("exit status 5"), expectedGitConfigErr: ErrTryingToUnsetAnOptionWhichDoesNotExist},
		{description: ErrTryingToUseAnInvalidRegexp.Error(), originalErr: errors.New("exit status 6"), expectedGitConfigErr: ErrTryingToUseAnInvalidRegexp},
	}

	for _, caseLoopVar := range cases {
		originalErr := caseLoopVar.originalErr
		gitCmd := caseLoopVar.gitCmd
		gitOutput := caseLoopVar.gitOutput
		expectedGitConfigErr := caseLoopVar.expectedGitConfigErr
		description := caseLoopVar.description

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			actualGitConfigErr := New(originalErr, gitCmd, gitOutput)

			if !reflect.DeepEqual(expectedGitConfigErr, actualGitConfigErr) {
				t.Errorf("expected: %s, got: %s", expectedGitConfigErr, actualGitConfigErr)
				t.Fail()
			}
		})
	}

}
