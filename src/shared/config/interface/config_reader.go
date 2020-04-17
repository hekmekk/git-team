package configinterface

import (
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
)

// Reader read the configuration
type Reader interface {
	Read() (config.Config, error)
}
