package config

import (
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	BaseDir          string
	TemplateFileName string
	StatusFileName   string
}

func Load() (Config, error) {
	baseDir, err := homedir.Expand("~/.config/git-team")
	if err != nil {
		return Config{}, err
	}

	return Config{BaseDir: baseDir, TemplateFileName: "COMMIT_TEMPLATE", StatusFileName: "status.toml"}, nil
}
