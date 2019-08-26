package enable

// Aborted aborted because no coauthors have been provided
type Aborted struct{}

// Succeeded successfully enabled git-team
type Succeeded struct{}

// Failed failed to enable git-team with Reason
type Failed struct {
	Reason []error
}
