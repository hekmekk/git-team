package writer

import (
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
)

// Writer write a single property
type Writer interface {
	SetActivationScope(scope activationscope.ActivationScope) error
}
