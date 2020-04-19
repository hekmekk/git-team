package commandadapter

// TODO: this should live somewhere else...

import (
	"fmt"

	"github.com/hekmekk/git-team/src/core/gitconfig"
)

// ResolveAliases convenience function to resolve multiple aliases and accumulate errors
func ResolveAliases(aliases []string) ([]string, []error) {
	return resolveAliases(ResolveAlias)(aliases)
}

func resolveAliases(resolveAlias func(string) (string, error)) func([]string) ([]string, []error) {
	return func(aliases []string) ([]string, []error) {
		var resolvedAliases []string
		var resolveErrors []error

		for _, alias := range aliases {
			var resolvedCoauthor, err = resolveAlias(alias)
			if err != nil {
				resolveErrors = append(resolveErrors, err)
			} else {
				resolvedAliases = append(resolvedAliases, resolvedCoauthor)
			}
		}

		return resolvedAliases, resolveErrors
	}
}

// ResolveAlias lookup "team.alias.<alias>" globally
func ResolveAlias(alias string) (string, error) {
	return resolveAlias(gitconfig.Get)(alias)
}

func resolveAlias(gitconfigGet func(string) (string, error)) func(string) (string, error) {
	return func(alias string) (string, error) {
		coauthor, err := gitconfigGet(fmt.Sprintf("team.alias.%s", alias))
		if err != nil || coauthor == "" {
			return "", fmt.Errorf("Failed to resolve alias team.alias.%s", alias)
		}

		return coauthor, nil
	}
}
