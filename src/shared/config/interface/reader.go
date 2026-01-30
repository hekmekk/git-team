package configinterface

import (
	config "github.com/hekmekk/git-team/v2/src/shared/config/entity/config"
)

// Reader read the configuration
type Reader interface {
	Read() (config.Config, error)
}
