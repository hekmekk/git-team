package policy

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

import (
	"github.com/hekmekk/git-team/v2/src/core/events"
)

// Policy the behavior that is applied when a command is issued
type Policy interface {
	Apply() events.Event
}
