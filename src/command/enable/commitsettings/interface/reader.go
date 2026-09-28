package commitsettingsinterface

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

import (
	"github.com/hekmekk/git-team/v2/src/command/enable/commitsettings/entity"
)

// Reader reads the internal commit settings
type Reader interface {
	Read() entity.CommitSettings
}
