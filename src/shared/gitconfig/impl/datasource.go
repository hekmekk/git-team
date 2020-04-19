package gitconfigimpl

import (
	gitconfiglegacy "github.com/hekmekk/git-team/src/shared/gitconfig/impl/legacy"
)

// DataSource read data directly from gitconfig
type DataSource struct{}

// NewDataSource construct new DataSource
func NewDataSource() DataSource {
	return DataSource{}
}

// Get read a setting
func (ds DataSource) Get(key string) (string, error) {
	return gitconfiglegacy.Get(key)
}
