package status

import (
	"github.com/hekmekk/git-team/src/core/state"
)

// StateRetrievalSucceeded successfully got the current state
type StateRetrievalSucceeded struct {
	State state.State
}

// StateRetrievalFailed failed to get the current state
type StateRetrievalFailed struct {
	Reason error
}
