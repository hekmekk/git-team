package state

import (
	"testing"
)

func TestIsEnabled(t *testing.T) {

	teamStatusToIsEnabled := make(map[teamStatus]bool)
	teamStatusToIsEnabled[enabled] = true
	teamStatusToIsEnabled[disabled] = false

	for status, expectedIsEnabled := range teamStatusToIsEnabled {
		isEnabled := State{Status: status}.IsEnabled()

		if isEnabled != expectedIsEnabled {
			t.Errorf("expected: %t, got: %t", expectedIsEnabled, isEnabled)
			t.Fail()
		}
	}
}
