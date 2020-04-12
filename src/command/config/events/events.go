package events

import (
	"github.com/hekmekk/git-team/src/core/config"
)

// RetrievalSucceeded successfully got the current state
type RetrievalSucceeded struct {
	Config config.Config
}

// RetrievalFailed failed to get the current state
type RetrievalFailed struct {
	Reason error
}
