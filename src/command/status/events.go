package status

import (
	state "github.com/hekmekk/git-team/src/shared/state/entity"
)

// StateRetrievalSucceeded successfully got the current state
type StateRetrievalSucceeded struct {
	State       state.State
	StateAsJson bool
}

// StateRetrievalFailed failed to get the current state
type StateRetrievalFailed struct {
	Reason error
}
