package handler

type AliasAdded struct {
	Alias    string
	CoAuthor string
}

// TODO: Refactor to match enable command pattern (don't use effects type tho, we have only one...)
func RunAddCommand(add func(string, string) error) func(string, string) (AliasAdded, error) {
	return func(alias string, coauthor string) (AliasAdded, error) {
		addErr := add(alias, coauthor)
		if addErr != nil {
			return AliasAdded{}, addErr
		}
		return AliasAdded{Alias: alias, CoAuthor: coauthor}, nil
	}
}
