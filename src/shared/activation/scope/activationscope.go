package entity

// Scope the scope of git team
type Scope int

const (
	// Global git team will be enabled and disabled globally
	Global Scope = iota
	// RepoLocal git team will be enabled and disabled for the current repository
	RepoLocal
	// Unknown no idea what to do with this value
	Unknown
)

func (scope Scope) String() string {
	names := [...]string{
		"global",
		"repo-local"}

	if scope < Global || scope > RepoLocal {
		return "unknown"
	}

	return names[scope]
}

// FromString factory method for Scope
func FromString(candidate string) Scope {
	switch candidate {
	case "global":
		return Global
	case "repo-local":
		return RepoLocal
	default:
		return Unknown
	}
}
