package datasink

import (
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
)

// GitconfigDataSink writes configuration to gitconfig
type GitconfigDataSink struct {
	GitConfigWriter gitconfig.Writer
}

// NewGitconfigDataSink constructs new GitconfigDataSink
func NewGitconfigDataSink(gitConfigWriter gitconfig.Writer) GitconfigDataSink {
	return GitconfigDataSink{GitConfigWriter: gitConfigWriter}
}

// SetActivationScope write activation-scope setting to gitconfig
func (ds GitconfigDataSink) SetActivationScope(scope activationscope.ActivationScope) error {
	return ds.GitConfigWriter.ReplaceAll("team.config.activation-scope", scope.String())
}
