package handler

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/hekmekk/git-team/core/git"
)

type AliasAdded struct {
	Alias    string
	CoAuthor string
}

func AddCommand(alias, coauthor string) (AliasAdded, error, int) {
	_, resolveErr := git.ResolveAlias(alias)
	if resolveErr == nil {
		return AliasAdded{}, errors.New(color.YellowString(fmt.Sprintf("Alias '%s' has already been added.", alias))), 0
	}

	addErr := git.AddAlias(alias, coauthor)
	if addErr != nil {
		return AliasAdded{}, addErr, -1
	}
	return AliasAdded{Alias: alias, CoAuthor: coauthor}, nil, 0
}
