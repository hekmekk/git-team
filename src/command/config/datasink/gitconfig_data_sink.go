package datasink

import (
	activationscope "github.com/hekmekk/git-team/src/command/config/entity/activationscope"
)

// GitconfigDataSink writes configuration to gitconfig
type GitconfigDataSink struct {
}

// NewGitconfigDataSink constructs new GitconfigDataSink
func NewGitconfigDataSink() GitconfigDataSink {
	return newGitconfigDataSink()
}

// for tests
func newGitconfigDataSink() GitconfigDataSink {
	return GitconfigDataSink{}
}

// SetActivationScope write activation-scope setting to gitconfig
func (ds GitconfigDataSink) SetActivationScope(scope activationscope.ActivationScope) error {
	return nil
}
