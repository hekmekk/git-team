package add

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

// AssignmentFailed assignment of coauthor to alias failed with Reason
type AssignmentFailed struct {
	Reason error
}

// AssignmentSucceeded assignment of coauthor to alias succeeded
type AssignmentSucceeded struct {
	Alias    string
	Coauthor string
}

// AssignmentAborted nothing changed
type AssignmentAborted struct {
}
