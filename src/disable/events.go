package disable

// Succeeded assignment of coauthor to alias succeeded
type Succeeded struct{}

// Failed failed to disable git-team
type Failed struct {
	Reason error
}
