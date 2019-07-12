package remove

// Command add a <Coauthor> under "team.alias.<Alias>"
type Command struct {
	Alias string
}

// Dependencies the real-world dependencies of the ExecutorFactory
type Dependencies struct {
	GitResolveAlias func(string) (string, error)
	GitRemoveAlias  func(string) error
}

// ExecutorFactory provisions a Command Processor
func ExecutorFactory(deps Dependencies) func(Command) error {
	return func(cmd Command) error {

		_, resolveErr := deps.GitResolveAlias(cmd.Alias)
		if resolveErr != nil {
			return nil
		}

		err := deps.GitRemoveAlias(cmd.Alias)
		if err != nil {
			return err
		}

		return nil
	}
}
