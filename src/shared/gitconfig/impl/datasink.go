package gitconfigimpl

import (
	gitconfiglegacy "github.com/hekmekk/git-team/src/shared/gitconfig/impl/legacy"
)

// DataSink write data directly to gitconfig
type DataSink struct{}

// NewDataSink construct new DataSink
func NewDataSink() DataSink {
	return DataSink{}
}

// ReplaceAll modify a setting
func (ds DataSink) ReplaceAll(key string, value string) error {
	return gitconfiglegacy.ReplaceAll(key, value)
}

// UnsetAll remove a setting
func (ds DataSink) UnsetAll(key string) error {
	return gitconfiglegacy.UnsetAll(key)
}
