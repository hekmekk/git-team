package remove

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

// DeAllocationFailed trying to remove an alias -> co-author assignment failed with Reason
type DeAllocationFailed struct {
	Reason error
}

// DeAllocationSucceeded successfully removed an alias -> co-author assignment
type DeAllocationSucceeded struct {
	Alias string
}
