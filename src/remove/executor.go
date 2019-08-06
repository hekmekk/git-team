package remove

import (
	git "github.com/hekmekk/git-team/src/gitconfig"
	giterror "github.com/hekmekk/git-team/src/gitconfig/error"
)

// Command remove an alias -> coauthor assignment
type Command struct {
	Alias string
}

type dependencies struct {
	GitRemoveAlias func(string) error
}

// Exec remove an alias -> coauthor assignment
func Exec(cmd Command) error {
	deps := dependencies{
		GitRemoveAlias: git.RemoveAlias,
	}
	return executorFactory(deps)(cmd)
}

func executorFactory(deps dependencies) func(Command) error {
	return func(cmd Command) error {
		err := deps.GitRemoveAlias(cmd.Alias)
		if err != nil && err.Error() != giterror.UnsetOptionWhichDoesNotExist {
			return err
		}

		return nil
	}
}
