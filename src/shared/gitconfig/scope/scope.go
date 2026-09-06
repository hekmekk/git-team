package scope

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
