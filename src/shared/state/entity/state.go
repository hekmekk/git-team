package stateentity

type teamStatus string

const (
	enabled  teamStatus = "enabled"
	disabled teamStatus = "disabled"
)

// State the state of git-team
type State struct {
	Status            teamStatus `json:"status"`
	Coauthors         []string   `json:"coAuthors"`
	PreviousHooksPath string     `json:"previousHooksPath"`
}

// NewStateEnabled the constructor for the enabled state
func NewStateEnabled(coauthors []string, previousHooksPath string) State {
	return State{Status: enabled, Coauthors: coauthors, PreviousHooksPath: previousHooksPath}
}

// NewStateDisabled the constructor for the disabled state
func NewStateDisabled() State {
	return State{Status: disabled, Coauthors: []string{}, PreviousHooksPath: ""}
}

// IsEnabled returns true if git-team is enabled
func (state State) IsEnabled() bool {
	return state.Status == enabled
}
