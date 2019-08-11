package config

import (
	"github.com/mitchellh/go-homedir"
)

// Config currently static config for git-team
type Config struct {
	BaseDir          string
	GitHooksPath     string
	TemplateFileName string
	StatusFileName   string
}

// Load loads the configuration file
func Load() (Config, error) {
	return executorFactory(dependencies{expandHomedir: homedir.Expand})()
}

type dependencies struct {
	expandHomedir func(string) (string, error)
}

func executorFactory(deps dependencies) func() (Config, error) {
	return func() (Config, error) {
		baseDir, err := deps.expandHomedir("~/.config/git-team")
		if err != nil {
			return Config{}, err
		}

		return Config{
			BaseDir:          baseDir,
			GitHooksPath:     "/usr/local/share/.config/git-team/hooks",
			TemplateFileName: "COMMIT_TEMPLATE",
			StatusFileName:   "status.toml",
		}, nil
	}
}
