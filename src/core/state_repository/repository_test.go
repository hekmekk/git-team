package staterepository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/core/state"
)

func TestQueryDisabled(t *testing.T) {
	deps := queryDependencies{
		gitconfigGet:    func(string) (string, error) { return "disabled", nil },
		gitconfigGetAll: func(string) ([]string, error) { return []string{}, nil },
	}

	expectedState := state.State{Status: "disabled", Coauthors: []string{}}

	state, err := query(deps)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, state) {
		t.Errorf("expected: %s, got: %s", expectedState, state)
		t.Fail()
	}
}

func TestQueryEnabled(t *testing.T) {
	activeCoauthors := []string{"Mr. Noujz <noujz@mr.se>"}

	deps := queryDependencies{
		gitconfigGet:    func(string) (string, error) { return "enabled", nil },
		gitconfigGetAll: func(string) ([]string, error) { return activeCoauthors, nil },
	}

	expectedState := state.State{Status: "enabled", Coauthors: activeCoauthors}

	state, err := query(deps)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, state) {
		t.Errorf("expected: %s, got: %s", expectedState, state)
		t.Fail()
	}
}

func TestQueryDisabledWhenStatusUnset(t *testing.T) {
	deps := queryDependencies{
		gitconfigGet:    func(string) (string, error) { return "", nil },
		gitconfigGetAll: func(string) ([]string, error) { return []string{}, nil },
	}

	expectedState := state.State{Status: "disabled", Coauthors: []string{}}

	state, err := query(deps)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, state) {
		t.Errorf("expected: %s, got: %s", expectedState, state)
		t.Fail()
	}
}

func TestPersistSucceeds(t *testing.T) {
	deps := persistDependencies{
		gitconfigUnsetAll:   func(string) error { return nil },
		gitconfigAdd:        func(string, string) error { return nil },
		gitconfigReplaceAll: func(string, string) error { return nil },
	}

	err := persist(deps, state.NewStateEnabled([]string{"CO-AUTHOR"}))

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestPersistSucceedsWhenTryingToRemoveNonExistingActiveCoauthorsFromGitConfig(t *testing.T) {
	deps := persistDependencies{
		gitconfigUnsetAll:   func(string) error { return errors.New("exit status 5") },
		gitconfigAdd:        func(string, string) error { return nil },
		gitconfigReplaceAll: func(string, string) error { return nil },
	}

	err := persist(deps, state.NewStateEnabled([]string{"CO-AUTHOR"}))

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestPersistFailsDueToAnotherUnsetAllFailure(t *testing.T) {
	expectedErr := errors.New("unset-all failure")

	deps := persistDependencies{
		gitconfigUnsetAll: func(string) error { return expectedErr },
	}

	err := persist(deps, state.NewStateEnabled([]string{"CO-AUTHOR"}))

	if expectedErr != err {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}
}

func TestPersistFailsDueToAddFailure(t *testing.T) {
	expectedErr := errors.New("add failure")

	deps := persistDependencies{
		gitconfigUnsetAll: func(string) error { return nil },
		gitconfigAdd:      func(string, string) error { return expectedErr },
	}

	err := persist(deps, state.NewStateEnabled([]string{"CO-AUTHOR"}))

	if expectedErr != err {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}
}

func TestPersistFailsDueReplaceAllFailure(t *testing.T) {
	expectedErr := errors.New("replace-all failure")

	deps := persistDependencies{
		gitconfigUnsetAll:   func(string) error { return nil },
		gitconfigAdd:        func(string, string) error { return nil },
		gitconfigReplaceAll: func(string, string) error { return expectedErr },
	}

	err := persist(deps, state.NewStateDisabled())

	if expectedErr != err {
		t.Errorf("expected: %s, got: %s", expectedErr, err)
		t.Fail()
	}
}
