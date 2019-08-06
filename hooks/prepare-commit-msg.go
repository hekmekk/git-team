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

type commitMsgSourceT string

const (
	message  commitMsgSourceT = "message"
	template commitMsgSourceT = "template"
	merge    commitMsgSourceT = "merge"
	commit   commitMsgSourceT = "commit"
	none     commitMsgSourceT = "none"
)

func main() {
	status, err := status.Fetch()
	if err != nil {
		panic(err)
	}

	if !status.IsEnabled() {
		os.Exit(0)
	}

	commitMsgSource, commitTemplate := parseArgs()

	switch commitMsgSource {
	case message:
		err := appendCoauthorsToCommitTemplate(commitTemplate, status.Coauthors)
		if err != nil {
			panic(err)
		}
	}

	os.Exit(0)
}

func parseArgs() (commitMsgSourceT, string) {
	flag.Parse()
	args := flag.Args()

	commitTemplate := args[0]
	commitMsgSource := none
	if len(args) >= 2 {
		commitMsgSource = commitMsgSourceT(args[1])
	}

	return commitMsgSource, commitTemplate
}

func appendCoauthorsToCommitTemplate(commitTemplate string, coauthors []string) error {
	f, err := os.OpenFile(commitTemplate, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(enableutils.PrepareForCommitMessage(coauthors)); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
