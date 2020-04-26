package stateinterface

import (
	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
)

// Writer persist the current state
type Writer interface {
	PersistEnabled(scope activationscope.ActivationScope, coauthors []string) error
	PersistDisabled(scope activationscope.ActivationScope) error
}
