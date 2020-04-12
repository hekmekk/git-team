package entity

import (
	activationscope "github.com/hekmekk/git-team/src/command/config/entity/activationscope"
)

// Config config for git-team
type Config struct {
	ActivationScope activationscope.ActivationScope
}
