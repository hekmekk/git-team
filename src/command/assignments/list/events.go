package list

import (
	"github.com/hekmekk/git-team/src/core/assignment"
)

// RetrievalFailed listing the available assignments failed
type RetrievalFailed struct {
	Reason error
}

// RetrievalSucceeded listing the available assignments succeeded
type RetrievalSucceeded struct {
	Assignments []assignment.Assignment
	OnlyAlias   bool
}
