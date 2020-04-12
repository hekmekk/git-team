package reader

import (
	config "github.com/hekmekk/git-team/src/command/config/entity/config"
)

// ConfigReader read the configuration
type ConfigReader interface {
	Read() (config.Config, error)
}
