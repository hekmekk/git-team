package gitconfigimpl

import (
	gitconfiglegacy "github.com/hekmekk/git-team/src/shared/gitconfig/impl/legacy"
	scope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

// DataSource read data directly from gitconfig
type DataSource struct{}

// NewDataSource construct new DataSource
func NewDataSource() DataSource {
	return DataSource{}
}

// Get read the first value for a key
func (ds DataSource) Get(scope scope.Scope, key string) (string, error) {
	return gitconfiglegacy.Get(scope, key)
}

// GetAll read all values for a key
func (ds DataSource) GetAll(scope scope.Scope, key string) ([]string, error) {
	return gitconfiglegacy.GetAll(scope, key)
}

// List show the entire config
func (ds DataSource) List(scope scope.Scope) (map[string]string, error) {
	return gitconfiglegacy.List(scope)
}
