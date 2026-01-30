package completion

import (
	"sort"
	"strings"

	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// AliasShellCompletion generate completion
type AliasShellCompletion struct {
	GitConfigReader gitconfig.Reader
}

// NewAliasShellCompletion construct new CoAuthorShellCompletion
func NewAliasShellCompletion(gitconfigReader gitconfig.Reader) AliasShellCompletion {
	return AliasShellCompletion{
		GitConfigReader: gitconfigReader,
	}
}

// Complete return not yet selected aliases
func (completion AliasShellCompletion) Complete(selectedAliases []string) []string {
	allAssignments, err := completion.GitConfigReader.GetRegexp(gitconfigscope.Global, "team.alias")

	if err != nil {
		return []string{}
	}

	remainingAliases := []string{}

	for rawAlias := range allAssignments {
		alias := strings.TrimPrefix(rawAlias, "team.alias.")
		isSelected := false
		for _, selectedAlias := range selectedAliases {
			if selectedAlias == alias {
				isSelected = true
			}
		}

		if !isSelected {
			remainingAliases = append(remainingAliases, alias)
		}
	}

	sort.Strings(remainingAliases)

	return remainingAliases
}
