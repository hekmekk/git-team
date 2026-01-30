package commitsettingsinterface

import (
	"github.com/hekmekk/git-team/v2/src/command/enable/commitsettings/entity"
)

// Reader reads the internal commit settings
type Reader interface {
	Read() entity.CommitSettings
}
