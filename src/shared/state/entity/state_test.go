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

import (
	"testing"
)

func TestIsEnabledShouldBeTrue(t *testing.T) {
	expectedIsEnabled := true
	isEnabled := NewStateEnabled([]string{}, "").IsEnabled()

	if expectedIsEnabled != isEnabled {
		t.Errorf("expected: %t, got: %t", expectedIsEnabled, isEnabled)
		t.Fail()
	}
}

func TestIsEnabledShouldBeFalse(t *testing.T) {
	expectedIsEnabled := false
	isEnabled := NewStateDisabled().IsEnabled()

	if expectedIsEnabled != isEnabled {
		t.Errorf("expected: %t, got: %t", expectedIsEnabled, isEnabled)
		t.Fail()
	}
}
