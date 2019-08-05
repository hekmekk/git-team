package main

import (
	"flag"
	"fmt"
	"os"
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
- due to the appending, the message will always be edited... No way to abort the commit...

Some samples:
-------------
git commit --amend: [DEBUG] args[0]=.git/COMMIT_EDITMSG
                    [DEBUG] args[1]=commit
                    [DEBUG] args[2]=HEAD[master 71dcf63] foo
git commit -m"foo": [DEBUG] args[0]=.git/COMMIT_EDITMSG
                    [DEBUG] args[1]=message[master 37985fe] foo
git commit -a     : [DEBUG] args[0]=.git/COMMIT_EDITMSG
git commit        : [DEBUG] args[0]=.git/COMMIT_EDITMSG

*/

const (
	MESSAGE  = "message"
	TEMPLATE = "template"
	MERGE    = "merge"
	COMMIT   = "commit"
	NONE     = "none"
)

func main() {

	flag.Parse()
	args := flag.Args()

	for index, arg := range flag.Args() {
		fmt.Printf("[DEBUG] args[%d]=%s\n", index, arg)
	}

	commitTemplate := args[0]
	commitMsgSource := NONE
	if len(args) >= 2 {
		commitMsgSource = args[1]

	}

	switch commitMsgSource {
	case NONE, MESSAGE, TEMPLATE:
		f, err := os.OpenFile(commitTemplate, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		if _, err := f.WriteString("\n\nCo-authored-by: No-one <no-one@foo.bar>"); err != nil {
			panic(err)
		}

		if err := f.Close(); err != nil {
			panic(err)
		}

	default:
		os.Exit(0)
	}
}
