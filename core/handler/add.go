package handler

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekmekk/git-team/core/git"
)

func AddCommand(alias, coauthor *string) {
	checkErr := git.SanityCheckCoauthor(*coauthor)
	if checkErr != nil {
		ToStderrAndExit(checkErr)
	}
	addErr := git.AddAlias(*alias, *coauthor)
	if addErr != nil {
		ToStderrAndExit(addErr)
	}
	color.Green(fmt.Sprintf("Alias '%s' -> %s has been added.", *alias, *coauthor))
}
