package status

import (
	"github.com/hekmekk/git-team/src/core/state"
)

// StateRetrievalSucceeded successfully queried the state repository for the current state
type StateRetrievalSucceeded struct {
	State state.State
}

// StateRetrievalFailed failed to query the state repository
type StateRetrievalFailed struct {
	Reason error
}
