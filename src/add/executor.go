package add

import (
	"github.com/hekmekk/git-team/src/validation"
)

// Args the arguments for the Executor returned by the ExecutorFactory
type Args struct {
	Alias    *string
	Coauthor *string
}

// Dependencies the real-world dependencies of the ExecutorFactory
type Dependencies struct {
	AddGitAlias func(string, string) error
}

// ExecutorFactory add a <Coauthor> under "team.alias.<Alias>"
func ExecutorFactory(deps Dependencies) func(Args) error {
	return func(args Args) error {
		checkErr := validation.SanityCheckCoauthor(*args.Coauthor)
		if checkErr != nil {
			return checkErr
		}

		addErr := deps.AddGitAlias(*args.Alias, *args.Coauthor)
		if addErr != nil {
			return addErr
		}
		return nil
	}
}
