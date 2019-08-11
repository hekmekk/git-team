package status

import (
	"testing"
)

func TestIsEnabled(t *testing.T) {

	teamStatusToIsEnabled := make(map[teamStatus]bool)
	teamStatusToIsEnabled[enabled] = true
	teamStatusToIsEnabled[disabled] = false

	for status, expectedIsEnabled := range teamStatusToIsEnabled {
		isEnabled := state{Status: status}.IsEnabled()

		if isEnabled != expectedIsEnabled {
			t.Errorf("expected: %t, got: %t", expectedIsEnabled, isEnabled)
			t.Fail()
		}
	}
}

func TestToStringEnabled(t *testing.T) {
	expectedString := "git-team enabled.\n\nCo-authors:\n-----------\nMr. Noujz <noujz@mr.se>\nMrs. Noujz <noujz@mrs.se>"

	str := state{Status: enabled, Coauthors: []string{"Mrs. Noujz <noujz@mrs.se>", "Mr. Noujz <noujz@mr.se>"}}.ToString()

	if str != expectedString {
		t.Errorf("expected: %s, got: %s", expectedString, str)
		t.Fail()
	}
}

func TestToStringDisabled(t *testing.T) {
	expectedString := "git-team disabled."

	str := state{Status: disabled}.ToString()

	if str != expectedString {
		t.Errorf("expected: %s, got: %s", expectedString, str)
		t.Fail()
	}
}
