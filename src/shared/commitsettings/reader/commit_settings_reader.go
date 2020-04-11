package reader

import (
	"github.com/hekmekk/git-team/src/shared/commitsettings/entity"
)

// CommitSettingsReader reads the internal commit settings
type CommitSettingsReader interface {
	Read() entity.CommitSettings
}
