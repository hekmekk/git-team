package main

import (
	"flag"
	"os"

	enableutils "github.com/hekmekk/git-team/src/command/enable/utils"
	config "github.com/hekmekk/git-team/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
	state "github.com/hekmekk/git-team/src/shared/state/impl"
)

type commitMsgSourceT string

const (
	commit   commitMsgSourceT = "commit"
	merge    commitMsgSourceT = "merge"
	message  commitMsgSourceT = "message"
	none     commitMsgSourceT = "none"
	squash   commitMsgSourceT = "squash"
	template commitMsgSourceT = "template"
)

func main() {
	configReader := config.NewGitconfigDataSource(gitconfig.NewDataSource())
	stateReader := state.NewGitConfigDataSource(gitconfig.NewDataSource())

	cfg, cfgReadErr := configReader.Read()
	if cfgReadErr != nil {
		panic(cfgReadErr)
	}

	status, statusQueryErr := stateReader.Query(cfg.ActivationScope)
	if statusQueryErr != nil {
		panic(statusQueryErr)
	}

	if !status.IsEnabled() {
		os.Exit(0)
	}

	commitMsgSource, commitTemplate := parseArgs()

	switch commitMsgSource {
	case message, merge, squash:
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
