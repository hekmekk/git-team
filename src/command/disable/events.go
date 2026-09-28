package disable

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

// Succeeded successfully disabled git-team
type Succeeded struct{}

// Failed failed to disable git-team
type Failed struct {
	Reason error
}
