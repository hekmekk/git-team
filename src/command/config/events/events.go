package events

import (
	config "github.com/hekmekk/git-team/src/command/config/entity/config"
)

// RetrievalSucceeded successfully got the current state
type RetrievalSucceeded struct {
	Config config.Config
}

// RetrievalFailed failed to get the current state
type RetrievalFailed struct {
	Reason error
}
