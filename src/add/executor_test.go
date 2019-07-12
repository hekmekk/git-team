package add

import (
	"errors"
	"testing"
)

func TestAddSucceeds(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	add := func(alias, coAuthor string) error {
		return nil
	}

	execAdd := ExecutorFactory(Dependencies{AddGitAlias: add})
	err := execAdd(Command{Alias: alias, Coauthor: coAuthor})

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestAddFailsBecauseUnderlyingGitAddFails(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	add := func(alias, coAuthor string) error {
		return errors.New("git add failed")
	}

	execAdd := ExecutorFactory(Dependencies{AddGitAlias: add})
	err := execAdd(Command{Alias: alias, Coauthor: coAuthor})

	if err == nil {
		t.Fail()
	}
}
