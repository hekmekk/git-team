package reader

import (
	"github.com/hekmekk/git-team/src/shared/internalconfig/entity"
)

// Reader reads the internal, static configuration
type Reader interface {
	Read() entity.InternalConfig
}
