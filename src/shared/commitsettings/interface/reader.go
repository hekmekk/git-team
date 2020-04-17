package commitsettingsinterface

import (
	"github.com/hekmekk/git-team/src/shared/commitsettings/entity"
)

// Reader reads the internal commit settings
type Reader interface {
	Read() entity.CommitSettings
}
