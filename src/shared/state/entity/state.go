package stateentity

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

type teamStatus string

const (
	enabled  teamStatus = "enabled"
	disabled teamStatus = "disabled"
)

// State the state of git-team
type State struct {
	Status            teamStatus
	Coauthors         []string
	PreviousHooksPath string
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
