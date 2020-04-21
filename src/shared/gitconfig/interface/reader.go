package gitconfiginterface

import (
	scope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

// Reader read git configuration settings
type Reader interface {
	Get(scope scope.Scope, key string) (string, error)
	GetAll(scope scope.Scope, key string) ([]string, error)
}
