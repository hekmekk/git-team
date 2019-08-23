package gitconfig

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestShouldExecuteGitConfigWithTheExpectedCommandLineArguments(t *testing.T) {
	staticArgs := []string{"config", "--null", "--global"}
	providedArgs := []string{"--one", "--two"}
	expectedArgs := append(staticArgs, providedArgs...)

	executor := func(args ...string) ([]byte, error) {
		if !reflect.DeepEqual(expectedArgs, args) {
			t.Errorf("expected: %s, received: %s", expectedArgs, args)
			t.Fail()
		}

		return nil, nil
	}

	execGitConfigFactory(executor)(providedArgs...)
}

func TestShouldReturnTheTwoLines(t *testing.T) {
	expectedLines := []string{"line1\n", "line2\n"}
	out := []byte(fmt.Sprintf("%s\000%s\000", expectedLines[0], expectedLines[1]))

	executor := func(args ...string) ([]byte, error) {
		return out, nil
	}

	lines, err := execGitConfigFactory(executor)("")
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

	lines, err := execGitConfigFactory(executor)("")
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

	lines, err := execGitConfigFactory(executor)("")
	if err == nil {
		t.Error("expected an error")
		t.Fail()
	}

	if lines != nil {
		t.Errorf("expected %s, received: %s", expectedLines, lines)
		t.Fail()
	}
}
