package gitconfig

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	scope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

func TestShouldExecuteGitConfigWithTheExpectedCommandLineArguments(t *testing.T) {
	t.Parallel()

	cases := []struct {
		scope             scope.Scope
		expectedScopeFlag string
	}{
		{scope.Global, "--global"},
		{scope.Local, "--local"},
	}

	for _, caseLoopVar := range cases {
		scope := caseLoopVar.scope
		expectedScopeFlag := caseLoopVar.expectedScopeFlag

		t.Run(scope.String(), func(t *testing.T) {
			providedOptions := []string{"--one", "--two"}
			expectedOptions := append([]string{expectedScopeFlag}, providedOptions...)

			executor := func(args ...string) ([]byte, error) {
				if !reflect.DeepEqual(expectedOptions, args) {
					t.Errorf("expected: %s, received: %s", expectedOptions, args)
					t.Fail()
				}

				return nil, nil
			}

			execGitConfigFactory(executor)(scope, providedOptions...)
		})
	}
}

func TestShouldReturnTheTwoLines(t *testing.T) {
	expectedLines := []string{"line1", "line2"}
	out := []byte(fmt.Sprintf("%s\n%s\n", expectedLines[0], expectedLines[1]))

	executor := func(args ...string) ([]byte, error) {
		return out, nil
	}

	lines, err := execGitConfigFactory(executor)(scope.Global, "")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedLines, lines) {
		t.Errorf("expected %s, received: %s", expectedLines, lines)
		t.Fail()
	}
}

func TestShouldReturnNoLines(t *testing.T) {
	expectedLines := []string{}

	executor := func(args ...string) ([]byte, error) {
		return []byte(""), nil
	}

	lines, err := execGitConfigFactory(executor)(scope.Global, "")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedLines, lines) {
		t.Errorf("expected %s, received: %s", expectedLines, lines)
		t.Fail()
	}
}

func TestShouldReturnAnError(t *testing.T) {
	expectedLines := []string{}

	executor := func(args ...string) ([]byte, error) {
		return nil, errors.New("executing the git config command failed")
	}

	lines, err := execGitConfigFactory(executor)(scope.Global, "")
	if err == nil {
		t.Error("expected an error")
		t.Fail()
	}

	if lines != nil {
		t.Errorf("expected %s, received: %s", expectedLines, lines)
		t.Fail()
	}
}
