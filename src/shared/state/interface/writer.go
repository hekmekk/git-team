package stateinterface

import (
	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
)

// Writer persist the current state
type Writer interface {
	PersistEnabled(scope activationscope.Scope, coauthors []string, previousHooksPath string) error
	PersistDisabled(scope activationscope.Scope) error
}
