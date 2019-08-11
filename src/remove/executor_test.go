package remove

import (
	"errors"
	"testing"
)

func TestRmSucceeds(t *testing.T) {
	alias := "mr"

	remove := func(alias string) error {
		return nil
	}

	execRemove := executorFactory(dependencies{GitRemoveAlias: remove})
	err := execRemove(Command{Alias: alias})

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestRmSucceedsWhenTryingToRemoveANonExistingAlias(t *testing.T) {
	alias := "mr"

	remove := func(alias string) error {
		return errors.New("exit status 5")
	}

	execRemove := executorFactory(dependencies{GitRemoveAlias: remove})
	err := execRemove(Command{Alias: alias})

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestRmFailsBecauseUnderlyingGitRemoveFails(t *testing.T) {
	alias := "mr"

	remove := func(alias string) error {
		return errors.New("git remove command failed")
	}

	execRemove := executorFactory(dependencies{GitRemoveAlias: remove})
	err := execRemove(Command{Alias: alias})

	if err == nil {
		t.Error("Expected remove to fail")
		t.Fail()
	}
}
