package handler

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hekmekk/git-team/core/git"
)

func AddCommand(alias, coauthor *string) {
	_, resolveErr := git.ResolveAlias(*alias)
	if resolveErr == nil {
		color.Yellow(fmt.Sprintf("Alias '%s' has already been added.", *alias))
		os.Exit(0)
	}

	checkErr := git.SanityCheckCoauthor(*coauthor)
	if checkErr != nil {
		ToStderrAndExit(checkErr)
	}
	addErr := git.AddAlias(*alias, *coauthor)
	if addErr != nil {
		ToStderrAndExit(addErr)
	}
	color.Green(fmt.Sprintf("Alias '%s' -> '%s' has been added.", *alias, *coauthor))
}
