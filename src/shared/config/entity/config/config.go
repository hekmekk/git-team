package entity

import (
	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
)

// Config config for git-team
type Config struct {
	ActivationScope activationscope.Scope
}
