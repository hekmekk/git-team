package main

import (
	"flag"
	"os"

	"github.com/hekmekk/git-team/src/enable/utils"
	"github.com/hekmekk/git-team/src/status"
)

/*
FRICKELAGE:
-----------
git config --global core.hooksPath ~/.config/git-team/hooks
mkdir ~/.config/git-team/hooks
go build -o ~/.config/git-team/hooks/prepare-commit-msg

git config --global --unset core.hooksPath

HOOK Docs:
----------
https://git-scm.com/docs/githooks#_prepare_commit_msg

ISSUE/THINGS TO CONSIDER:
-------------------------
- be careful not to overwrite core.hooksPath if it's set already
- be careful not to install prepare-commit-msg if it's present
*/

const (
	message  = "message"
	template = "template"
	merge    = "merge"
	commit   = "commit"
	none     = "none"
)

func main() {
	status, err := status.Fetch()
	if err != nil {
		panic(err)
	}

	if !status.IsEnabled() {
		os.Exit(0)
	}

	flag.Parse()
	args := flag.Args()

	commitTemplate := args[0]
	commitMsgSource := none
	if len(args) >= 2 {
		commitMsgSource = args[1]
	}

	switch commitMsgSource {
	case message:
		f, err := os.OpenFile(commitTemplate, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		if _, err := f.WriteString(enableutils.PrepareForCommitMessage(status.Coauthors)); err != nil {
			panic(err)
		}

		if err := f.Close(); err != nil {
			panic(err)
		}

	default:
		os.Exit(0)
	}
}
