package state

type teamStatus string

const (
	enabled  teamStatus = "enabled"
	disabled teamStatus = "disabled"
)

// State the state of git-team
type State struct {
	Status    teamStatus `toml:"status"`
	Coauthors []string   `toml:"co-authors"`
}

// NewStateEnabled the constructor for the enabled state
func NewStateEnabled(coauthors []string) State {
	return State{Status: enabled, Coauthors: coauthors}
}

// NewStateDisabled the constructor for the disabled state
func NewStateDisabled() State {
	return State{Status: disabled, Coauthors: []string{}}
}

// IsEnabled returns true if git-team is enabled
func (state State) IsEnabled() bool {
	return state.Status == enabled
}
