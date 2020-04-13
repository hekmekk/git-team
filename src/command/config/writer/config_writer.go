package writer

import (
	activationscope "github.com/hekmekk/git-team/src/command/config/entity/activationscope"
)

// ConfigWriter write a single property
type ConfigWriter interface {
	SetActivationScope(scope activationscope.ActivationScope) error
}
