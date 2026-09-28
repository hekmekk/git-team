package configinterface

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

import (
	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
)

// Writer write a single property
type Writer interface {
	SetActivationScope(scope activationscope.Scope) error
}
