package reader

import (
	"github.com/hekmekk/git-team/src/shared/internalconfig/entity"
)

// ConfigReader reads the internal, static configuration
type ConfigReader interface {
	Read() entity.InternalConfig
}
