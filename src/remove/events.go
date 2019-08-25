package remove

// DeAllocationFailed trying to remove an alias -> co-author assignment failed with Reason
type DeAllocationFailed struct {
	Reason error
}

// DeAllocationSucceeded successfully removed an alias -> co-author assignment
type DeAllocationSucceeded struct {
	Alias string
}
