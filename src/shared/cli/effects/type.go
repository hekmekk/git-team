package effects

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

// Effect a side effect
type Effect interface {
	Run() error
}
