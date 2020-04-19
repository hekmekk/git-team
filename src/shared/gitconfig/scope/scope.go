package scope

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

// Flag translate the gitconfig scope to the respective flag
func (scope Scope) Flag() string {
	flags := [...]string{
		"--global",
		"--local"}

	return flags[scope]
}
