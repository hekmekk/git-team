package add

// Command add a <Coauthor> under "team.alias.<Alias>"
type Command struct {
	Alias    string
	Coauthor string
}

// Dependencies the real-world dependencies of the ExecutorFactory
type Dependencies struct {
	AddGitAlias func(string, string) error
}

// ExecutorFactory provisions a Command Processor
func ExecutorFactory(deps Dependencies) func(Command) error {
	return func(cmd Command) error {
		addErr := deps.AddGitAlias(cmd.Alias, cmd.Coauthor)
		if addErr != nil {
			return addErr
		}
		return nil
	}
}
