package add

import (
	"github.com/hekmekk/git-team/src/gitconfig"
	"github.com/hekmekk/git-team/src/validation"
)

// Args the arguments of the Runner
type Args struct {
	Alias    *string
	Coauthor *string
}

// Dependencies the dependencies of the Runner
type Dependencies struct {
	AddGitAlias func(string, string) error
}

// Definition the command, arguments, and dependencies
type Definition struct {
	CommandName string
	Args        Args
	Deps        Dependencies
}

// New the constructor for Definition
func New(name string, alias, coauthor *string) Definition {
	return Definition{
		CommandName: name,
		Args: Args{
			Alias:    alias,
			Coauthor: coauthor,
		},
		Deps: Dependencies{
			AddGitAlias: gitconfig.AddAlias,
		},
	}
}

// AssignmentFailed assignment of coauthor to alias failed with Reason
type AssignmentFailed struct {
	Reason error
}

// AssignmentSucceeded assignment of coauthor to alias succeeded
type AssignmentSucceeded struct {
	Alias    string
	Coauthor string
}

// Apply assign a co-author to an alias
func Apply(deps Dependencies, args Args) interface{} {
	alias := *args.Alias
	coauthor := *args.Coauthor
	err := assignCoauthorToAlias(deps, alias, coauthor)
	if err != nil {
		return AssignmentFailed{Reason: err}
	}

	return AssignmentSucceeded{Alias: alias, Coauthor: coauthor}
}

func assignCoauthorToAlias(deps Dependencies, alias string, coauthor string) error {
	checkErr := validation.SanityCheckCoauthor(coauthor) // TODO: Add as dependency. Has no side effects but doesn't need to be tested here.
	if checkErr != nil {
		return checkErr
	}

	addErr := deps.AddGitAlias(alias, coauthor)
	if addErr != nil {
		return addErr
	}
	return nil
}
