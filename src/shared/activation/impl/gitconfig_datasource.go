package activationvalidatorimpl

import (
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// GitConfigDataSource use gitconfig to determine activation validation properties
type GitConfigDataSource struct {
	GitConfigReader gitconfig.Reader
}

// NewGitConfigDataSource constructor for GitConfigDataSource
func NewGitConfigDataSource(gitConfigReader gitconfig.Reader) GitConfigDataSource {
	return GitConfigDataSource{GitConfigReader: gitConfigReader}
}

// IsInsideAGitRepository find out if the current working dir is a git repository
func (ds GitConfigDataSource) IsInsideAGitRepository() bool {
	_, err := ds.GitConfigReader.List(gitconfigscope.Local)
	return err == nil
}
