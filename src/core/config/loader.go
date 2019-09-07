package config

import (
	"fmt"
	"os"
)

// Config currently static config for git-team
type Config struct {
	BaseDir          string
	GitHooksPath     string
	TemplateFileName string
	StatusFileName   string
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
			BaseDir:          fmt.Sprintf("%s/.config/git-team", deps.getEnv("HOME")),
			GitHooksPath:     "/usr/local/etc/git-team/hooks",
			TemplateFileName: "COMMIT_TEMPLATE",
			StatusFileName:   "status.toml",
		}
	}
}
