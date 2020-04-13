package datasink

import (
	activationscope "github.com/hekmekk/git-team/src/command/config/entity/activationscope"
	coregitconfig "github.com/hekmekk/git-team/src/core/gitconfig"
)

// GitconfigDataSink writes configuration to gitconfig
type GitconfigDataSink struct {
	RawGitConfigWriter coregitconfig.RawWriter
}

// NewGitconfigDataSink constructs new GitconfigDataSink
func NewGitconfigDataSink() GitconfigDataSink {
	return newGitconfigDataSink(coregitconfig.NewRawDataSink())
}

// for tests
func newGitconfigDataSink(rawGitConfigWriter coregitconfig.RawWriter) GitconfigDataSink {
	return GitconfigDataSink{RawGitConfigWriter: rawGitConfigWriter}
}

// SetActivationScope write activation-scope setting to gitconfig
func (ds GitconfigDataSink) SetActivationScope(scope activationscope.ActivationScope) error {
	return ds.RawGitConfigWriter.Set("team.config.activation-scope", scope.String())
}
