package configinterface

import (
	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
)

// Writer write a single property
type Writer interface {
	SetActivationScope(scope activationscope.ActivationScope) error
}
