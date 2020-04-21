package stateentity

import (
	"testing"
)

func TestIsEnabledShouldBeTrue(t *testing.T) {
	expectedIsEnabled := true
	isEnabled := NewStateEnabled([]string{}).IsEnabled()

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
