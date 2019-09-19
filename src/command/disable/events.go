package disable

// Succeeded successfully disabled git-team
type Succeeded struct{}

// Failed failed to disable git-team
type Failed struct {
	Reason error
}
