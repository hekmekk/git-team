package handler

type AliasAdded struct {
	Alias    string
	CoAuthor string
}

func RunAddCommand(add func(string, string) error) func(string, string) (AliasAdded, error) {
	return func(alias string, coauthor string) (AliasAdded, error) {
		addErr := add(alias, coauthor)
		if addErr != nil {
			return AliasAdded{}, addErr
		}
		return AliasAdded{Alias: alias, CoAuthor: coauthor}, nil
	}
}
