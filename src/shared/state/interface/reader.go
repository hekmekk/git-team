package stateinterface

import (
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	state "github.com/hekmekk/git-team/src/shared/state/entity"
)

// Reader retrive the current state
type Reader interface {
	Query(scope activationscope.ActivationScope) (state.State, error)
}
