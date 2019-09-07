package staterepository

import (
	"errors"

	"github.com/hekmekk/git-team/src/core/gitconfig"
	giterror "github.com/hekmekk/git-team/src/core/gitconfig/error"
	"github.com/hekmekk/git-team/src/core/state"
)

// Query read the current state from file
func Query() (state.State, error) {
	deps := queryDependencies{
		gitconfigGet:    gitconfig.Get,
		gitconfigGetAll: gitconfig.GetAll,
	}
	return query(deps)
}

type queryDependencies struct {
	gitconfigGet    func(string) (string, error)
	gitconfigGetAll func(string) ([]string, error)
}

const (
	keyStatus          = "team.state.status"
	keyActiveCoauthors = "team.state.active-coauthors"
)

func query(deps queryDependencies) (state.State, error) {
	status, err := deps.gitconfigGet(keyStatus)
	if err != nil || "disabled" == status || "" == status {
		return state.NewStateDisabled(), nil
	}

	activeCoauthors, err := deps.gitconfigGetAll(keyActiveCoauthors)
	if err != nil {
		return state.State{}, errors.New("no active co-authors found")
	}

	return state.NewStateEnabled(activeCoauthors), nil
}

// PersistEnabled persist the current state as enabled
func PersistEnabled(coauthors []string) error {
	return Persist(state.NewStateEnabled(coauthors))
}

// PersistDisabled persist the current state as disabled
func PersistDisabled() error {
	return Persist(state.NewStateDisabled())
}

// Persist persist the current state
func Persist(state state.State) error {
	deps := persistDependencies{
		gitconfigAdd:        gitconfig.Add,
		gitconfigReplaceAll: gitconfig.ReplaceAll,
		gitconfigUnsetAll:   gitconfig.UnsetAll,
	}
	return persist(deps, state)
}

type persistDependencies struct {
	gitconfigAdd        func(key, value string) error
	gitconfigReplaceAll func(key, value string) error
	gitconfigUnsetAll   func(key string) error
}

func persist(deps persistDependencies, state state.State) error {
	if err := deps.gitconfigUnsetAll(keyActiveCoauthors); err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
		return err
	}

	for _, coauthor := range state.Coauthors {
		if err := deps.gitconfigAdd(keyActiveCoauthors, coauthor); err != nil {
			return err
		}
	}

	if err := deps.gitconfigReplaceAll(keyStatus, string(state.Status)); err != nil {
		return err
	}

	return nil
}
