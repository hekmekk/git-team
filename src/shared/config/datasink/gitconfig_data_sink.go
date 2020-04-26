package datasink

import (
	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
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
	return ds.GitConfigWriter.ReplaceAll(gitconfigscope.Global, "team.config.activation-scope", scope.String())
}
