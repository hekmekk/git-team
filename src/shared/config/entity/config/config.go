package entity

import (
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
)

// Config config for git-team
type Config struct {
	ActivationScope activationscope.ActivationScope
}
