package stateinterface

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

import (
	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	state "github.com/hekmekk/git-team/v2/src/shared/state/entity"
)

// Reader retrive the current state
type Reader interface {
	Query(scope activationscope.Scope) (state.State, error)
}
