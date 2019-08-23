package gitconfig

import (
	"errors"
	"testing"
)

func TestShouldReturnTheAssignedCoAuthor(t *testing.T) {
	mr := "mr"
	mrNoujz := "Mr. Noujz <noujz@mr.se>"

	execSucceeds := func(args ...string) ([]string, error) { return []string{mrNoujz}, nil }

	coauthor, err := resolveAlias(execSucceeds)(mr)

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

	execSucceeds := func(args ...string) ([]string, error) { return []string{}, nil }

	coauthor, err := resolveAlias(execSucceeds)(mr)

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

	execFails := func(args ...string) ([]string, error) { return nil, errors.New("git command failed") }

	coauthor, err := resolveAlias(execFails)(mr)

	if err == nil {
		t.Error("expected an error")
		t.Fail()
	}

	if coauthor != "" {
		t.Errorf("expected: %s, received: %s", "", coauthor)
		t.Fail()
	}
}
