package entity

// ActivationScope the scope of git team
type ActivationScope int

func (scope ActivationScope) String() string {
	names := [...]string{
		"global",
		"repo-local"}

	if scope < Global || scope > RepoLocal {
		return "unknown"
	}

	return names[scope]
}

const (
	// Global git team will be enabled and disabled globally
	Global ActivationScope = iota
	// RepoLocal git team will be enabled and disabled for the current repository
	RepoLocal
)
