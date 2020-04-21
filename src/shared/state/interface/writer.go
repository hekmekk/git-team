package stateinterface

import (
	scope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

// TODO: should accept ActivationScope

// Writer persist the current state
type Writer interface {
	PersistEnabled(scope scope.Scope, coauthors []string) error
	PersistDisabled(scope scope.Scope) error
}
