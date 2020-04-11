package config

// TODO: split this up... see the imports in the related files, need to put files into diff packages
// TODO: not only this.. check recent tests and all other files... weird imports and missing structure

import (
	"fmt"
	"os"
)

// Repository the repository for the git team config
type Repository interface {
	Query() (Config, error)
}

type StaticValueDataSource struct {
}

func NewStaticValueDataSource() StaticValueDataSource {
	return StaticValueDataSource{}
}

func (ds StaticValueDataSource) Query() (Config, error) {
	cfg := Load()
	return cfg, nil
}

// ReadOnlyProperties read only properties of the config
type ReadOnlyProperties struct {
	GitTeamCommitTemplatePath string
	GitTeamHooksPath          string
}

// ActivationScope the scope of git team
type ActivationScope int

func (scope ActivationScope) String() string {
	names := [...]string{
		"global",
		"repo-local"}

	if scope < Global || scope > RepoLocal {
		return "unknown"
	}

	return names[scope]
}

const (
	// Global git team will be enabled and disabled globally
	Global ActivationScope = iota
	// RepoLocal git team will be enabled and disabled for the current repository
	RepoLocal
)

// ReadWriteProperties read/write properties of the config
type ReadWriteProperties struct {
	ActivationScope ActivationScope
}

// Config config for git-team
type Config struct {
	Ro ReadOnlyProperties
	Rw ReadWriteProperties
}

// Load loads the configuration file
func Load() Config {
	return executorFactory(dependencies{getEnv: os.Getenv})()
}

type dependencies struct {
	getEnv func(string) string
}

func executorFactory(deps dependencies) func() Config {
	return func() Config {
		return Config{
			Ro: ReadOnlyProperties{
				GitTeamCommitTemplatePath: fmt.Sprintf("%s/.config/git-team/COMMIT_TEMPLATE", deps.getEnv("HOME")),
				GitTeamHooksPath:          "/usr/local/etc/git-team/hooks",
			},
			Rw: ReadWriteProperties{
				ActivationScope: Global,
			},
		}
	}
}
