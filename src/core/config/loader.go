package config

import (
	"fmt"
	"os"
)

// Config currently static config for git-team
type Config struct {
	GitTeamCommitTemplatePath string
	GitTeamHooksPath          string
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
			GitTeamCommitTemplatePath: fmt.Sprintf("%s/.config/git-team/COMMIT_TEMPLATE", deps.getEnv("HOME")),
			GitTeamHooksPath:          "/usr/local/etc/git-team/hooks",
		}
	}
}
