package gitconfigimpl

import (
	coregitconfig "github.com/hekmekk/git-team/src/core/gitconfig"
)

// DataSource read data directly from gitconfig
type DataSource struct{}

// NewDataSource construct new DataSource
func NewDataSource() DataSource {
	return DataSource{}
}

// Get read a setting
func (ds DataSource) Get(key string) (string, error) {
	return coregitconfig.Get(key)
}
