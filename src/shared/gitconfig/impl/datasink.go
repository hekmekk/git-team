package gitconfigimpl

import (
	gitconfiglegacy "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl/legacy"
	scope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// DataSink write data directly to gitconfig
type DataSink struct{}

// NewDataSink construct new DataSink
func NewDataSink() DataSink {
	return DataSink{}
}

// Add add a value to a key
func (ds DataSink) Add(scope scope.Scope, key string, value string) error {
	return gitconfiglegacy.Add(scope, key, value)
}

// ReplaceAll modify a setting
func (ds DataSink) ReplaceAll(scope scope.Scope, key string, value string) error {
	return gitconfiglegacy.ReplaceAll(scope, key, value)
}

// UnsetAll remove a setting
func (ds DataSink) UnsetAll(scope scope.Scope, key string) error {
	return gitconfiglegacy.UnsetAll(scope, key)
}
