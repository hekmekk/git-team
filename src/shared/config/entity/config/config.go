package entity

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

import (
	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
)

// Config config for git-team
type Config struct {
	ActivationScope activationscope.Scope
}
