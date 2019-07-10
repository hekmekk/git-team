package handler

type AddCommand struct {
	Alias    string
	Coauthor string
}

type AddEffect struct {
	AddGitAlias func(string, string) error
}

func ExecutorFactory(effect AddEffect) func(AddCommand) error {
	return func(cmd AddCommand) error {
		addErr := effect.AddGitAlias(cmd.Alias, cmd.Coauthor)
		if addErr != nil {
			return addErr
		}
		return nil
	}
}
