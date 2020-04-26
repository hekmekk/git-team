package entity

// TODO: should this live in its own module?
// TODO: rename to Scope -> activationscope.Scope

// ActivationScope the scope of git team
type ActivationScope int

const (
	// Global git team will be enabled and disabled globally
	Global ActivationScope = iota
	// RepoLocal git team will be enabled and disabled for the current repository
	RepoLocal
	// Unknown no idea what to do with this value
	Unknown
)

func (scope ActivationScope) String() string {
	names := [...]string{
		"global",
		"repo-local"}

	if scope < Global || scope > RepoLocal {
		return "unknown"
	}

	return names[scope]
}

// FromString factory method for ActivationScope
func FromString(candidate string) ActivationScope {
	switch candidate {
	case "global":
		return Global
	case "repo-local":
		return RepoLocal
	default:
		return Unknown
	}
}
