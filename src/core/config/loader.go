package config

import (
	"fmt"
	"os"
)

// ReadOnlyProperties read only properties of the config
type ReadOnlyProperties struct {
	GitTeamCommitTemplatePath string
	GitTeamHooksPath          string
}

// Config config for git-team
type Config struct {
	Ro ReadOnlyProperties
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
		}
	}
}
