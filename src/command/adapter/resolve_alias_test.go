package commandadapter

import (
	"errors"
	"testing"
)

func TestShouldReturnTheAssignedCoAuthor(t *testing.T) {
	mr := "mr"
	mrNoujz := "Mr. Noujz <noujz@mr.se>"

	gitconfigGet := func(args string) (string, error) { return mrNoujz, nil }

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

	gitconfigGet := func(alias string) (string, error) { return "", nil }

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

	gitconfigGet := func(alias string) (string, error) { return "", errors.New("git command failed") }

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
