package handler

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekmekk/git-team/core/git"
)

func RemoveCommand(alias *string) {
	err := git.RemoveAlias(*alias)
	if err != nil {
		ToStderrAndExit(err)
	}
	color.Red(fmt.Sprintf("Alias '%s' has been removed.", *alias))
}
