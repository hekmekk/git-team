package datasource

import (
	activationscope "github.com/hekmekk/git-team/src/command/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/command/config/entity/config"
)

// GitconfigDataSource reads configuration from git config
type GitconfigDataSource struct {
}

// NewGitconfigDataSource constructs new GitconfigDataSource
func NewGitconfigDataSource() GitconfigDataSource {
	return newGitconfigDataSource()
}

// for tests
func newGitconfigDataSource() GitconfigDataSource {
	return GitconfigDataSource{}
}

func (ds GitconfigDataSource) Read() (config.Config, error) {
	cfg := config.Config{
		ActivationScope: activationscope.Global,
	}

	return cfg, nil
}
