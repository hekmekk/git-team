package gitconfig

// Scope scope of gitconfig
type Scope int

const (
	// Global global gitconfig
	Global Scope = iota
	// Local local gitconfig
	Local
)

func (scope Scope) String() string {
	names := [...]string{
		"global",
		"local"}

	return names[scope]
}
