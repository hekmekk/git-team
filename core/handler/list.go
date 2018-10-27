package handler

import (
	"github.com/hekmekk/git-team/core/git"
)

func ListCommand() {
	err := git.ListAlias()
	if err != nil {
		ToStderrAndExit(err)
	}
}
