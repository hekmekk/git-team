package commandadapter

import (
	"errors"
	"testing"

	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

func TestShouldReturnTheAssignedCoAuthor(t *testing.T) {
	mr := "mr"
	mrNoujz := "Mr. Noujz <noujz@mr.se>"

	gitconfigGet := func(scope gitconfigscope.Scope, args string) (string, error) {
		if scope != gitconfigscope.Global {
			return "", errors.New("wrong scope")
		}
		return mrNoujz, nil
	}

	coauthor, err := resolveAlias(gitconfigGet)(mr)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if coauthor != mrNoujz {
		t.Errorf("expected: %s, received: %s", mrNoujz, coauthor)
		t.Fail()
	}
}

func TestShouldReturnErrorIfNoAssignmentIsFound(t *testing.T) {
	mr := "mr"

	gitconfigGet := func(scope gitconfigscope.Scope, alias string) (string, error) {
		if scope != gitconfigscope.Global {
			return "", errors.New("wrong scope")
		}
		return "", nil
	}

	coauthor, err := resolveAlias(gitconfigGet)(mr)

	if err == nil {
		t.Error("expected an error")
		t.Fail()
	}

	if coauthor != "" {
		t.Errorf("expected: %s, received: %s", "", coauthor)
		t.Fail()
	}
}

func TestShouldReturnErrorIfResolvingFails(t *testing.T) {
	mr := "mr"
	expectedErr := errors.New("Failed to resolve alias team.alias.mr")

	gitconfigGet := func(scope gitconfigscope.Scope, alias string) (string, error) {
		if scope != gitconfigscope.Global {
			return "", errors.New("wrong scope")
		}
		return "", errors.New("git command failed")
	}

	coauthor, err := resolveAlias(gitconfigGet)(mr)

	if err.Error() != expectedErr.Error() {
		t.Errorf("expected: %s, received: %s", expectedErr, err)
		t.Fail()
	}

	if coauthor != "" {
		t.Errorf("expected: %s, received: %s", "", coauthor)
		t.Fail()
	}
}
