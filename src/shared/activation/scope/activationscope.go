package entity

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
