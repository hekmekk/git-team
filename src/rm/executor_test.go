package rm

import (
	"errors"
	"testing"
)

func TestRmSucceeds(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	resolve := func(alias string) (string, error) {
		return coAuthor, nil
	}

	remove := func(alias string) error {
		return nil
	}

	execAdd := ExecutorFactory(RemoveDependencies{GitResolveAlias: resolve, GitRemoveAlias: remove})
	err := execAdd(RemoveCommand{Alias: alias})

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestRmSucceedsRegardlessOfResolveError(t *testing.T) {
	alias := "mr"

	resolve := func(alias string) (string, error) {
		return "", errors.New("no such alias")
	}

	remove := func(alias string) error {
		return nil
	}

	execAdd := ExecutorFactory(RemoveDependencies{GitResolveAlias: resolve, GitRemoveAlias: remove})
	err := execAdd(RemoveCommand{Alias: alias})

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestRmFailsBecauseUnderlyingGitRemoveFails(t *testing.T) {
	alias := "mr"
	coAuthor := "Mr. Noujz <noujz@mr.se>"

	resolve := func(alias string) (string, error) {
		return coAuthor, nil
	}

	remove := func(alias string) error {
		return errors.New("git remove command failed")
	}

	execAdd := ExecutorFactory(RemoveDependencies{GitResolveAlias: resolve, GitRemoveAlias: remove})
	err := execAdd(RemoveCommand{Alias: alias})

	if err == nil {
		t.Error("Expected remove to fail")
		t.Fail()
	}
}
