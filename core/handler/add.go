package handler

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

type AliasAdded struct {
	Alias    string
	CoAuthor string
}

func RunAddCommand(resolve func(string) (string, error), add func(string, string) error) func(string, string) (AliasAdded, error, int) {
	return func(alias string, coauthor string) (AliasAdded, error, int) {
		_, resolveErr := resolve(alias)
		if resolveErr == nil {
			return AliasAdded{}, errors.New(color.YellowString(fmt.Sprintf("Alias '%s' has already been added.", alias))), 0
		}

		addErr := add(alias, coauthor)
		if addErr != nil {
			return AliasAdded{}, addErr, -1
		}
		return AliasAdded{Alias: alias, CoAuthor: coauthor}, nil, 0
	}
}
