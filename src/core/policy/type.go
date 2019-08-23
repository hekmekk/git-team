package policy

import (
	"github.com/hekmekk/git-team/src/core/events"
)

// Policy the behavior that is applied when a command is issued
type Policy interface {
	Apply() events.Event
}
