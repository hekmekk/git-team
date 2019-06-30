package handler

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/hekmekk/git-team/core/git"
)

func ListCommand() {
	lines, err := git.ListAlias()
	if err != nil {
		lines = make([]string, 0)
	}

	blackBold := color.New(color.FgBlack).Add(color.Bold)
	blackBold.Println("Aliases:")
	blackBold.Println("--------")

	for _, v := range lines {
		aliasToCoauthor := strings.Split(strings.TrimRight(v, "\n"), "\n")
		color.Magenta(fmt.Sprintf("'%s' -> '%s'", strings.TrimPrefix(aliasToCoauthor[0], "team.alias."), aliasToCoauthor[1]))
	}
}
