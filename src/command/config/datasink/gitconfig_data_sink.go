package datasink

import (
	activationscope "github.com/hekmekk/git-team/src/command/config/entity/activationscope"
	coregitconfig "github.com/hekmekk/git-team/src/core/gitconfig"
)

// GitconfigDataSink writes configuration to gitconfig
type GitconfigDataSink struct {
	GitSettingsWriter coregitconfig.SettingsWriter
}

// NewGitconfigDataSink constructs new GitconfigDataSink
func NewGitconfigDataSink() GitconfigDataSink {
	return newGitconfigDataSink(coregitconfig.NewDataSink())
}

// for tests
func newGitconfigDataSink(gitSettingsWriter coregitconfig.SettingsWriter) GitconfigDataSink {
	return GitconfigDataSink{GitSettingsWriter: gitSettingsWriter}
}

// SetActivationScope write activation-scope setting to gitconfig
func (ds GitconfigDataSink) SetActivationScope(scope activationscope.ActivationScope) error {
	return ds.GitSettingsWriter.Set("team.config.activation-scope", scope.String())
}
