package commitsettingsinterface

import (
	"github.com/hekmekk/git-team/src/command/enable/commitsettings/entity"
)

// Reader reads the internal commit settings
type Reader interface {
	Read() entity.CommitSettings
}
