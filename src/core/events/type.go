package events

import (
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
)

// Event an event emitted when applying a Policy
type Event interface{}

// EventAdapter the boundary between policy (Event) and cli (Effect)
type EventAdapter interface {
	MapEventToEffects(Event) []effects.Effect
}
