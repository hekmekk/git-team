package activationvalidatorinterface

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

// Validator check activation validation properties
type Validator interface {
	IsInsideAGitRepository() bool
}
