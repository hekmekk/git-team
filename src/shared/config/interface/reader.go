package configinterface

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

import (
	config "github.com/hekmekk/git-team/v2/src/shared/config/entity/config"
)

// Reader read the configuration
type Reader interface {
	Read() (config.Config, error)
}
