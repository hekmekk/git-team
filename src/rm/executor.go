package rm

type RemoveCommand struct {
	Alias string
}

type RemoveDependencies struct {
	GitResolveAlias func(string) (string, error)
	GitRemoveAlias  func(string) error
}

func ExecutorFactory(deps RemoveDependencies) func(RemoveCommand) error {
	return func(cmd RemoveCommand) error {

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
